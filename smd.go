package main

import (
	"math"

	"github.com/aoisensi/src2gltf/flag"
	"github.com/aoisensi/src2gltf/smd"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func smd2gltf(smd *smd.SMD, out string) error {
	doc := gltf.NewDocument()
	doc.Scenes[0].Nodes = []uint32{0}
	doc.Nodes = []*gltf.Node{{Mesh: gltf.Index(0)}}
	doc.Meshes = []*gltf.Mesh{{
		Primitives: []*gltf.Primitive{
			{
				Attributes: gltf.Attribute{
					"POSITION": 1,
				},
				Indices: gltf.Index(0),
			},
		},
	}}
	mxf32 := float32(math.MaxFloat32)
	glvsc := uint16(0)                           //GL VetexeS Counter
	glvs := make([][3]float32, 0, 65536)         //GL VertexeS
	glvsm := make(map[[3]float32]uint16, 65536)  //GL VertexeS Map
	glts := make([]uint16, 0, 65536)             //GL TriangleS
	glmaxv := [3]float32{-mxf32, -mxf32, -mxf32} //GL Max Vertex
	glminv := [3]float32{mxf32, mxf32, mxf32}    //GL Min Vertex
	for _, ts := range smd.Triangles {
		for _, t := range ts.Vertexes {
			v3mulp(&t.Pos, float32(flag.Scale))
			for a := 0; a < 3; a++ {
				if glminv[a] > t.Pos[a] {
					glminv[a] = t.Pos[a]
				}
				if glmaxv[a] < t.Pos[a] {
					glmaxv[a] = t.Pos[a]
				}
			}
			i, ok := glvsm[t.Pos]
			if ok {
				glts = append(glts, i)
				continue
			}
			glvs = append(glvs, t.Pos)
			glvsm[t.Pos] = glvsc
			glts = append(glts, glvsc)
			glvsc++
		}
	}
	modeler.WriteIndices(doc, glts)
	modeler.WritePosition(doc, glvs)
	doc.Buffers[0].EmbeddedResource()
	doc.Accessors = []*gltf.Accessor{
		{
			BufferView:    gltf.Index(0),
			ComponentType: gltf.ComponentUshort,
			Count:         uint32(len(glts)),
			Type:          gltf.AccessorScalar,
			Min:           []float32{0},
			Max:           []float32{float32(len(glvs) - 1)},
		},
		{
			BufferView:    gltf.Index(1),
			ComponentType: gltf.ComponentFloat,
			Count:         uint32(len(glvs)),
			Type:          gltf.AccessorVec3,
			Min:           glminv[:],
			Max:           glmaxv[:],
		},
	}
	return gltf.Save(doc, out)
}

func v3mulp(a *[3]float32, b float32) {
	a[0] *= b
	a[1] *= b
	a[2] *= b
}
