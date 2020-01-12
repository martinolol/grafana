package plugins

import (
	"context"
	"encoding/json"
	"path"

	pluginModel "github.com/grafana/grafana-plugin-model/go/renderer"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/util/errutil"
)

type RendererPlugin struct {
	PluginBase

	Executable string `json:"executable,omitempty"`
	GrpcPlugin pluginModel.RendererPlugin
}

func (r *RendererPlugin) Load(decoder *json.Decoder, pluginDir string) error {
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	if err := r.registerPlugin(pluginDir); err != nil {
		return err
	}

	cmd := ComposePluginStartCommmand("plugin_start")
	fullpath := path.Join(r.PluginDir, cmd)
	descriptor := backendplugin.NewRendererPluginDescriptor(r.Id, fullpath, &backendplugin.PluginCallbacks{
		LegacyCallback: r.onLegacyPluginStart,
	})
	if err := backendplugin.Register(descriptor); err != nil {
		return errutil.Wrapf(err, "Failed to register backend plugin")
	}

	Renderer = r
	return nil
}

func (r *RendererPlugin) Start(ctx context.Context) error {
	if err := backendplugin.StartPlugin(ctx, r.Id); err != nil {
		return errutil.Wrapf(err, "Failed to start renderer plugin")
	}

	return nil
}

func (r *RendererPlugin) onLegacyPluginStart(pluginID string, client *backendplugin.LegacyClient, logger log.Logger) error {
	r.GrpcPlugin = client.RendererPlugin
	return nil
}
