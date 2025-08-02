package utils

import (
	"bytes"
	"io"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
)

func ParseYaml(namespace, yamlContent string) (unstructureList []*unstructured.Unstructured, err error) {
	d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yamlContent), 4096)
	var unstructureObj *unstructured.Unstructured
	for {
		unstructureObj, err = GetUnstructured(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		unstructureList = append(unstructureList, unstructureObj)
	}
	return unstructureList, nil
}

func GetUnstructured(d *uyaml.YAMLOrJSONDecoder) (unstructureObj *unstructured.Unstructured, err error) {
	var rawObj runtime.RawExtension
	err = d.Decode(&rawObj)
	if err == io.EOF {
		return
	}
	if err != nil {
		err = errorx.NewDefaultError("Decode is err: %v", err.Error())
		return
	}
	obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
	if err != nil {
		err = errorx.NewDefaultError("Rawobj is err: %v", err.Error())
		return
	}
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		err = errorx.NewDefaultError("Tounstructured is err %v", err.Error())
		return
	}
	unstructureObj = &unstructured.Unstructured{Object: unstructuredMap}
	return
}
