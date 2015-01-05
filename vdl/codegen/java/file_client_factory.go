package java

import (
	"bytes"
	"log"
	"path"

	"v.io/core/veyron2/vdl/compile"
)

const clientFactoryTmpl = `// This file was auto-generated by the veyron vdl tool.
// Source(s):  {{ .Sources }}
package {{ .PackagePath }};

/* Factory for binding to {{ .ServiceName }}Client interfaces. */
{{.AccessModifier}} final class {{ .ServiceName }}ClientFactory {
    public static {{ .ServiceName }}Client bind(final java.lang.String name) throws io.v.core.veyron2.VeyronException {
        return bind(name, null);
    }
    public static {{ .ServiceName }}Client bind(final java.lang.String name, final io.v.core.veyron2.Options veyronOpts) throws io.v.core.veyron2.VeyronException {
        io.v.core.veyron2.ipc.Client client = null;
        if (veyronOpts != null && veyronOpts.get(io.v.core.veyron2.OptionDefs.CLIENT) != null) {
            client = veyronOpts.get(io.v.core.veyron2.OptionDefs.CLIENT, io.v.core.veyron2.ipc.Client.class);
        } else {
            client = io.v.core.veyron2.VRuntime.getClient();
        }
        return new {{ .StubName }}(client, name);
    }
}
`

// genJavaClientFactoryFile generates the Java file containing client bindings for
// all interfaces in the provided package.
func genJavaClientFactoryFile(iface *compile.Interface, env *compile.Env) JavaFileInfo {
	javaServiceName := toUpperCamelCase(iface.Name)
	data := struct {
		AccessModifier string
		Sources        string
		ServiceName    string
		PackagePath    string
		StubName       string
	}{
		AccessModifier: accessModifierForName(iface.Name),
		Sources:        iface.File.BaseName,
		ServiceName:    javaServiceName,
		PackagePath:    javaPath(javaGenPkgPath(iface.File.Package.Path)),
		StubName:       javaPath(javaGenPkgPath(path.Join(iface.File.Package.Path, iface.Name+"ClientStub"))),
	}
	var buf bytes.Buffer
	err := parseTmpl("client factory", clientFactoryTmpl).Execute(&buf, data)
	if err != nil {
		log.Fatalf("vdl: couldn't execute client template: %v", err)
	}
	return JavaFileInfo{
		Name: javaServiceName + "ClientFactory.java",
		Data: buf.Bytes(),
	}
}
