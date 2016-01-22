package image

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/builder/dockerfile"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/docker/docker/runconfig"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	registrytypes "github.com/docker/engine-api/types/registry"
	"golang.org/x/net/context"
)

func (s *imageRouter) postCommit(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	if err := httputils.CheckForJSON(r); err != nil {
		return err
	}

	cname := r.Form.Get("container")

	pause := httputils.BoolValue(r, "pause")
	version := httputils.VersionFromContext(ctx)
	if r.FormValue("pause") == "" && version.GreaterThanOrEqualTo("1.13") {
		pause = true
	}

	c, _, _, err := runconfig.DecodeContainerConfig(r.Body)
	if err != nil && err != io.EOF { //Do not fail if body is empty.
		return err
	}
	if c == nil {
		c = &container.Config{}
	}

	if !s.daemon.Exists(cname) {
		return derr.ErrorCodeNoSuchContainer.WithArgs(cname)
	}

	newConfig, err := dockerfile.BuildFromConfig(c, r.Form["changes"])
	if err != nil {
		return err
	}

	commitCfg := &types.ContainerCommitConfig{
		Pause:        pause,
		Repo:         r.Form.Get("repo"),
		Tag:          r.Form.Get("tag"),
		Author:       r.Form.Get("author"),
		Comment:      r.Form.Get("comment"),
		Config:       newConfig,
		MergeConfigs: true,
	}

	imgID, err := s.daemon.Commit(cname, commitCfg)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusCreated, &types.ContainerCommitResponse{
		ID: string(imgID),
	})
}

// Creates an image from Pull or from Import
func (s *imageRouter) postImagesCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	var (
		image   = r.Form.Get("fromImage")
		repo    = r.Form.Get("repo")
		tag     = r.Form.Get("tag")
		message = r.Form.Get("message")
		err     error
		output  = ioutils.NewWriteFlusher(w)
	)
	defer output.Close()

	w.Header().Set("Content-Type", "application/json")

	if image != "" { //pull
		// Special case: "pull -a" may send an image name with a
		// trailing :. This is ugly, but let's not break API
		// compatibility.
		image = strings.TrimSuffix(image, ":")

		var ref reference.Named
		ref, err = reference.ParseNamed(image)
		if err == nil {
			if tag != "" {
				// The "tag" could actually be a digest.
				var dgst digest.Digest
				dgst, err = digest.ParseDigest(tag)
				if err == nil {
					ref, err = reference.WithDigest(ref, dgst)
				} else {
					ref, err = reference.WithTag(ref, tag)
				}
			}
			if err == nil {
				metaHeaders := map[string][]string{}
				for k, v := range r.Header {
					if strings.HasPrefix(k, "X-Meta-") {
						metaHeaders[k] = v
					}
				}

				authConfigs := make(map[string]types.AuthConfig)
				authConfigs, err = s.getAuthConfigs(ref, r, false, false, "")
				if err != nil {
					return err
				}

				err = s.daemon.PullImage(ref, metaHeaders, authConfigs, output)
			}
		}
		// Check the error from pulling an image to make sure the request
		// was authorized. Modify the status if the request was
		// unauthorized to respond with 401 rather than 500.
		if err != nil && isAuthorizedError(err) {
			err = errcode.ErrorCodeUnauthorized.WithMessage(fmt.Sprintf("Authentication is required: %s", err))
		}
	} else { //import
		var newRef reference.Named
		if repo != "" {
			var err error
			newRef, err = reference.ParseNamed(repo)
			if err != nil {
				return err
			}

			if _, isCanonical := newRef.(reference.Canonical); isCanonical {
				return errors.New("cannot import digest reference")
			}

			if tag != "" {
				newRef, err = reference.WithTag(newRef, tag)
				if err != nil {
					return err
				}
			}
		}

		src := r.Form.Get("fromSrc")

		// 'err' MUST NOT be defined within this block, we need any error
		// generated from the download to be available to the output
		// stream processing below
		var newConfig *container.Config
		newConfig, err = dockerfile.BuildFromConfig(&container.Config{}, r.Form["changes"])
		if err != nil {
			return err
		}

		err = s.daemon.ImportImage(src, newRef, message, r.Body, output, newConfig)
	}
	if err != nil {
		if !output.Flushed() {
			return err
		}
		sf := streamformatter.NewJSONStreamFormatter()
		output.Write(sf.FormatError(err))
	}

	return nil
}

func (s *imageRouter) postImagesPush(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	metaHeaders := map[string][]string{}
	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-Meta-") {
			metaHeaders[k] = v
		}
	}
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	ref, err := reference.ParseNamed(vars["name"])
	if err != nil {
		return err
	}
	tag := r.Form.Get("tag")
	if tag != "" {
		// Push by digest is not supported, so only tags are supported.
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			return err
		}
	}

	authConfigs, err := s.getAuthConfigs(ref, r, true, false, "")
	if err != nil {
		return err
	}

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()

	w.Header().Set("Content-Type", "application/json")

	if err := s.daemon.PushImage(ref, metaHeaders, authConfigs, output); err != nil {
		if !output.Flushed() {
			return err
		}
		sf := streamformatter.NewJSONStreamFormatter()
		output.Write(sf.FormatError(err))
	}
	return nil
}

func (s *imageRouter) getImagesGet(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/x-tar")

	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	var names []string
	if name, ok := vars["name"]; ok {
		names = []string{name}
	} else {
		names = r.Form["names"]
	}

	if err := s.daemon.ExportImage(names, output); err != nil {
		if !output.Flushed() {
			return err
		}
		sf := streamformatter.NewJSONStreamFormatter()
		output.Write(sf.FormatError(err))
	}
	return nil
}

func (s *imageRouter) postImagesLoad(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	quiet := httputils.BoolValueOrDefault(r, "quiet", true)
	w.Header().Set("Content-Type", "application/json")
	return s.daemon.LoadImage(r.Body, w, quiet)
}

func (s *imageRouter) deleteImages(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	name := vars["name"]

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("image name cannot be blank")
	}

	force := httputils.BoolValue(r, "force")
	prune := !httputils.BoolValue(r, "noprune")

	list, err := s.daemon.ImageDelete(name, force, prune)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, list)
}

func (s *imageRouter) getImagesByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	imageInspect, err := s.daemon.LookupImage(vars["name"])
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, imageInspect)
}

func (s *imageRouter) getImagesJSON(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	// FIXME: The filter parameter could just be a match filter
	images, err := s.daemon.Images(r.Form.Get("filters"), r.Form.Get("filter"), httputils.BoolValue(r, "all"))
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, images)
}

func (s *imageRouter) getImagesHistory(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	name := vars["name"]
	history, err := s.daemon.ImageHistory(name)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, history)
}

func (s *imageRouter) postImagesTag(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	repo := r.Form.Get("repo")
	tag := r.Form.Get("tag")
	newTag, err := reference.WithName(repo)
	if err != nil {
		return err
	}
	if tag != "" {
		if newTag, err = reference.WithTag(newTag, tag); err != nil {
			return err
		}
	}
	if err := s.daemon.TagImage(newTag, vars["name"]); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *imageRouter) getImagesSearch(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	var (
		headers = map[string][]string{}
		term    = r.Form.Get("term")
	)

	authConfigs, err := s.getAuthConfigs(nil, r, false, true, term)
	if err != nil {
		return err
	}

	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-Meta-") {
			headers[k] = v
		}
	}
	results, err := s.daemon.SearchRegistryForImages(term, authConfigs, headers, httputils.BoolValue(r, "noIndex"))
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, results)
}

func isAuthorizedError(err error) bool {
	if urlError, ok := err.(*url.Error); ok {
		err = urlError.Err
	}

	if dError, ok := err.(errcode.Error); ok {
		if dError.ErrorCode() == errcode.ErrorCodeUnauthorized {
			return true
		}
	}
	return false
}

func (s *imageRouter) getAuthConfigs(ref reference.Named, r *http.Request, backward, search bool, searchTerm string) (map[string]types.AuthConfig, error) {
	authEncoded := r.Header.Get("X-Registry-Auth")
	authConfigs := make(map[string]types.AuthConfig)
	if authEncoded != "" {
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(&authConfigs); err != nil {
			// for a pull it is not an error if no auth was given
			// to increase compatibility with the existing api it is defaulting to be empty
			authConfigs = make(map[string]types.AuthConfig)
		}
	} else if backward {
		// the old format is supported for compatibility if there was no authConfig header
		if err := json.NewDecoder(r.Body).Decode(&authConfigs); err != nil {
			return nil, fmt.Errorf("Bad parameters and missing X-Registry-Auth: %v", err)
		}
	}
	// maybe client just sends one auth config
	// try to resolve just one auth config...
	authConfig := types.AuthConfig{}
	if len(authConfigs) == 0 {
		if authEncoded != "" {
			authJSONSingle := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
			if err := json.NewDecoder(authJSONSingle).Decode(&authConfig); err != nil {
				// for a pull it is not an error if no auth was given
				// to increase compatibility with the existing api it is defaulting to be empty
				authConfig = types.AuthConfig{}
			}
		} else if backward {
			// the old format is supported for compatibility if there was no authConfig header
			if err := json.NewDecoder(r.Body).Decode(&authConfig); err != nil {
				return nil, fmt.Errorf("Bad parameters and missing X-Registry-Auth: %v", err)
			}
		}
	}

	if len(authConfigs) == 0 {
		var (
			indexInfo *registrytypes.IndexInfo
			err       error
		)
		if search {
			indexInfo, err = registry.ParseSearchIndexInfo(searchTerm)
			if err != nil {
				return nil, err
			}
		} else {
			repoInfo, err := registry.ParseRepositoryInfo(ref)
			if err != nil {
				return nil, err
			}
			indexInfo = repoInfo.Index
		}

		// search default to nil if no X-Registry-Auth
		if authEncoded == "" && search {
			// XXX(runcom): this should be `return nil, nil` cause session v1
			// looks for single authConfig nilness
			return authConfigs, nil
		}
		// this is the case when we fully qualify images
		// XXX(runcom): add a test to ensure I can `docker push docker.io/runcom/something`
		// otherwise using indexInfo.Name in map as a key do not work for docker.io
		authConfigs[registry.GetAuthConfigKey(indexInfo)] = authConfig
	}

	return authConfigs, nil
}
