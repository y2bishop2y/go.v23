package java

import (
	"bytes"
	"log"
	"path"

	"veyron.io/veyron/veyron2/vdl/compile"
	"veyron.io/veyron/veyron2/vdl/vdlutil"
)

const serverWrapperTmpl = `// This file was auto-generated by the veyron vdl tool.
// Source(s):  {{ .Source }}
package {{ .PackagePath }};

{{ .AccessModifier }} final class {{ .ServiceName }}ServerWrapper {

    private final {{ .FullServiceName }}Server server;

{{/* Define fields to hold each of the embedded server wrappers*/}}
{{ range $embed := .Embeds }}
    {{/* e.g. private final com.somepackage.gen_impl.ArithStub stubArith; */}}
    private final {{ $embed.WrapperClassName }} {{ $embed.LocalWrapperVarName }};
    {{ end }}

    public {{ .ServiceName }}ServerWrapper(final {{ .FullServiceName }}Server server) {
        this.server = server;
        {{/* Initialize the embeded server wrappers */}}
        {{ range $embed := .Embeds }}
        this.{{ $embed.LocalWrapperVarName }} = new {{ $embed.WrapperClassName }}(server);
        {{ end }}
    }

    /**
     * Returns a description of this server.
     */
     public io.veyron.veyron.veyron2.ipc.ServiceSignature signature(io.veyron.veyron.veyron2.ipc.ServerCall call) throws io.veyron.veyron.veyron2.VeyronException {
         throw new io.veyron.veyron.veyron2.VeyronException("Signature method not yet supported for Java servers");
     }

    /**
     * Returns all tags associated with the provided method or null if the method isn't implemented
     * by this server.
     */
    public java.lang.Object[] getMethodTags(final io.veyron.veyron.veyron2.ipc.ServerCall call, final java.lang.String method) {
        {{ range $methodName, $tags := .MethodTags }}
        if ("{{ $methodName }}".equals(method)) {
            return new java.lang.Object[] {
                {{ range $tag := $tags }} {{ $tag }}, {{ end }}
            };
        }
        {{ end }}
        {{ range $embed := .Embeds }}
        {
            final java.lang.Object[] tags = this.{{ $embed.LocalWrapperVarName }}.getMethodTags(call, method);
            if (tags != null) {
                return tags;
            }
        }
        {{ end }}
        return null;  // method not found
    }

     {{/* Iterate over methods defined directly in the body of this server */}}
    {{ range $method := .Methods }}
    {{ $method.AccessModifier }} {{ $method.RetType }} {{ $method.Name }}(final io.veyron.veyron.veyron2.ipc.ServerCall call{{ $method.DeclarationArgs }}) throws io.veyron.veyron.veyron2.VeyronException {
        {{ if $method.IsStreaming }}
        final io.veyron.veyron.veyron2.vdl.Stream<{{ $method.SendType }}, {{ $method.RecvType }}> _stream = new io.veyron.veyron.veyron2.vdl.Stream<{{ $method.SendType }}, {{ $method.RecvType }}>() {
            @Override
            public void send({{ $method.SendType }} item) throws io.veyron.veyron.veyron2.VeyronException {
                final java.lang.reflect.Type type = new com.google.common.reflect.TypeToken< {{ $method.SendType }} >() {}.getType();
                call.send(item, type);
            }
            @Override
            public {{ $method.RecvType }} recv() throws java.io.EOFException, io.veyron.veyron.veyron2.VeyronException {
                final java.lang.reflect.Type type = new com.google.common.reflect.TypeToken< {{ $method.RecvType }} >() {}.getType();
                final java.lang.Object result = call.recv(type);
                try {
                    return ({{ $method.RecvType }})result;
                } catch (java.lang.ClassCastException e) {
                    throw new io.veyron.veyron.veyron2.VeyronException("Unexpected result type: " + result.getClass().getCanonicalName());
                }
            }
        };
        {{ end }} {{/* end if $method.IsStreaming */}}
        {{ if $method.Returns }} return {{ end }} this.server.{{ $method.Name }}( call {{ $method.CallingArgs }} {{ if $method.IsStreaming }} ,_stream {{ end }} );
    }
{{end}}

{{/* Iterate over methods from embeded servers and generate code to delegate the work */}}
{{ range $eMethod := .EmbedMethods }}
    {{ $eMethod.AccessModifier }} {{ $eMethod.RetType }} {{ $eMethod.Name }}(final io.veyron.veyron.veyron2.ipc.ServerCall call{{ $eMethod.DeclarationArgs }}) throws io.veyron.veyron.veyron2.VeyronException {
        {{/* e.g. return this.stubArith.cosine(call, [args], options) */}}
        {{ if $eMethod.Returns }}return{{ end }}  this.{{ $eMethod.LocalWrapperVarName }}.{{ $eMethod.Name }}(call{{ $eMethod.CallingArgs }});
    }
{{ end }} {{/* end range .EmbedMethods */}}

}
`

type serverWrapperMethod struct {
	AccessModifier  string
	CallingArgs     string
	DeclarationArgs string
	IsStreaming     bool
	Name            string
	RecvType        string
	RetType         string
	Returns         bool
	SendType        string
}

type serverWrapperEmbedMethod struct {
	AccessModifier      string
	CallingArgs         string
	DeclarationArgs     string
	LocalWrapperVarName string
	Name                string
	RetType             string
	Returns             bool
}

type serverWrapperEmbed struct {
	LocalWrapperVarName string
	WrapperClassName    string
}

func processServerWrapperMethod(iface *compile.Interface, method *compile.Method, env *compile.Env) serverWrapperMethod {
	return serverWrapperMethod{
		AccessModifier:  accessModifierForName(method.Name),
		CallingArgs:     javaCallingArgStr(method.InArgs, true),
		DeclarationArgs: javaDeclarationArgStr(method.InArgs, env, true),
		IsStreaming:     isStreamingMethod(method),
		Name:            vdlutil.ToCamelCase(method.Name),
		RecvType:        javaType(method.OutStream, true, env),
		RetType:         clientInterfaceOutArg(iface, method, true, env),
		Returns:         len(method.OutArgs) >= 2,
		SendType:        javaType(method.InStream, true, env),
	}
}

func processServerWrapperEmbedMethod(iface *compile.Interface, embedMethod *compile.Method, env *compile.Env) serverWrapperEmbedMethod {
	return serverWrapperEmbedMethod{
		AccessModifier:      accessModifierForName(embedMethod.Name),
		CallingArgs:         javaCallingArgStr(embedMethod.InArgs, true),
		DeclarationArgs:     javaDeclarationArgStr(embedMethod.InArgs, env, true),
		LocalWrapperVarName: vdlutil.ToCamelCase(iface.Name) + "Wrapper",
		Name:                vdlutil.ToCamelCase(embedMethod.Name),
		RetType:             clientInterfaceOutArg(iface, embedMethod, true, env),
		Returns:             len(embedMethod.OutArgs) >= 2,
	}
}

// genJavaServerWrapperFile generates a java file containing a server wrapper for the specified
// interface.
func genJavaServerWrapperFile(iface *compile.Interface, env *compile.Env) JavaFileInfo {
	embeds := []serverWrapperEmbed{}
	for _, embed := range allEmbeddedIfaces(iface) {
		embeds = append(embeds, serverWrapperEmbed{
			WrapperClassName:    javaPath(javaGenPkgPath(path.Join(embed.File.Package.Path, toUpperCamelCase(embed.Name+"ServerWrapper")))),
			LocalWrapperVarName: vdlutil.ToCamelCase(embed.Name) + "Wrapper",
		})
	}
	methodTags := make(map[string][]string)
	// Add generated methods to the tag map:
	methodTags["signature"] = []string{}
	methodTags["getMethodTags"] = []string{}
	// Copy method tags off of the interface.
	for _, method := range iface.Methods {
		tags := make([]string, len(method.Tags))
		for i, tagVal := range method.Tags {
			tags[i] = javaConstVal(tagVal, env)
		}
		methodTags[vdlutil.ToCamelCase(method.Name)] = tags
	}
	embedMethods := []serverWrapperEmbedMethod{}
	for _, embedMao := range dedupedEmbeddedMethodAndOrigins(iface) {
		embedMethods = append(embedMethods, processServerWrapperEmbedMethod(embedMao.Origin, embedMao.Method, env))
	}
	methods := make([]serverWrapperMethod, len(iface.Methods))
	for i, method := range iface.Methods {
		methods[i] = processServerWrapperMethod(iface, method, env)
	}
	javaServiceName := toUpperCamelCase(iface.Name)
	data := struct {
		AccessModifier  string
		EmbedMethods    []serverWrapperEmbedMethod
		Embeds          []serverWrapperEmbed
		FullServiceName string
		Methods         []serverWrapperMethod
		MethodTags      map[string][]string
		PackagePath     string
		ServiceName     string
		Source          string
	}{
		AccessModifier:  accessModifierForName(iface.Name),
		EmbedMethods:    embedMethods,
		Embeds:          embeds,
		FullServiceName: javaPath(interfaceFullyQualifiedName(iface)),
		Methods:         methods,
		MethodTags:      methodTags,
		PackagePath:     javaPath(javaGenPkgPath(iface.File.Package.Path)),
		ServiceName:     javaServiceName,
		Source:          iface.File.BaseName,
	}
	var buf bytes.Buffer
	err := parseTmpl("server wrapper", serverWrapperTmpl).Execute(&buf, data)
	if err != nil {
		log.Fatalf("vdl: couldn't execute server wrapper template: %v", err)
	}
	return JavaFileInfo{
		Name: javaServiceName + "ServerWrapper.java",
		Data: buf.Bytes(),
	}
}
