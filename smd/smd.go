package smd

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type SMD struct {
	Nodes     []*Node
	Skeleton  Skeleton
	Triangles []*Triangle
}

type decoder struct {
	sc  *bufio.Scanner
	smd *SMD
}

func Decode(r io.Reader) (*SMD, error) {
	sc := bufio.NewScanner(r)
	d := &decoder{sc: sc, smd: new(SMD)}
	line, err := d.Scan()
	if err != nil {
		return nil, err
	}
	if line != "version 1" {
		return nil, errors.New("not vaild smd file or not supported version")
	}
	for {
		line, err = d.Scan()
		if err != nil {
			if err == io.EOF {
				return d.smd, nil
			}
			return nil, err
		}
		switch line {
		case "nodes":
			d.ReadNodes()
		case "skeleton":
			d.ReadSkeleton()
		case "triangles":
			d.ReadTriangles()
		}
	}
}

func (d *decoder) ReadNodes() error {
	d.smd.Nodes = make([]*Node, 0, 1024)
	for {
		line, err := d.Scan()
		if err != nil {
			return err
		}
		if line == "end" {
			return nil
		}
		vv := strings.Split(line, " ")
		if len(vv) != 3 {
			return errors.New("failed load nodes")
		}
		id, err := strconv.Atoi(vv[0])
		if err != nil {
			return err
		}
		pid, err := strconv.Atoi(vv[2])
		if err != nil {
			return err
		}
		node := &Node{
			ID:       id,
			Name:     strings.Trim(vv[1], "\""),
			ParentID: pid,
		}
		d.smd.Nodes = append(d.smd.Nodes, node)
	}
}

func (d *decoder) ReadSkeleton() error {
	d.smd.Skeleton = make(map[int][]*SkeletonBone, 1024)
	time := -1000
	for {
		line, err := d.Scan()
		if err != nil {
			return err
		}
		if line == "end" {
			return nil
		}
		vv := strings.Split(line, " ")
		switch len(vv) {
		case 2:
			if vv[0] != "time" {
				return errors.New("failed load skeleton")
			}
			time, err = strconv.Atoi(vv[1])
			if err != nil {
				return err
			}
			d.smd.Skeleton[time] = make([]*SkeletonBone, 0, 256)
		case 7:
			bid, err := strconv.Atoi(vv[0])
			if err != nil {
				return err
			}
			pos, err := ssto3f(vv[1:])
			if err != nil {
				return err
			}
			rot, err := ssto3f(vv[4:])
			if err != nil {
				return err
			}
			bone := &SkeletonBone{
				BoneID: bid,
				Pos:    pos,
				Rot:    rot,
			}
			d.smd.Skeleton[time] = append(d.smd.Skeleton[time], bone)
		default:
			return errors.New("failed load skeleton")
		}
	}
}

func (d *decoder) ReadTriangles() error {
	d.smd.Triangles = make([]*Triangle, 0, 1048576)
	for {
		tri := &Triangle{}
		line, err := d.Scan()
		if err != nil {
			return err
		}
		if line == "end" {
			return nil
		}
		tri.Material = line
		for i := 0; i < 3; i++ {
			tv := new(TriangleVertex)
			line, err := d.Scan()
			if err != nil {
				return err
			}
			vv := strings.Split(line, " ")
			if len(vv) < 9 {
				return errors.New("failed load triangles")
			}
			tv.ParentBoneID, err = strconv.Atoi(vv[0])
			if err != nil {
				return err
			}
			tv.Pos, err = ssto3f(vv[1:])
			if err != nil {
				return err
			}
			tv.Norm, err = ssto3f(vv[4:])
			if err != nil {
				return err
			}
			tv.UV, err = ssto2f(vv[7:])
			if err != nil {
				return err
			}
			if len(vv) == 9 {
				tri.Vertexes[i] = tv
				continue
			}
			ls, err := strconv.Atoi(vv[9])
			links := make([]*TriangleVertexLink, ls)
			for j := 0; j < ls; j++ {
				bid, err := strconv.Atoi(vv[10+j*2])
				if err != nil {
					return err
				}
				w, err := strconv.ParseFloat(vv[11+j*2], 32)
				links[j] = &TriangleVertexLink{
					BoneID: bid,
					Weight: float32(w),
				}
			}
			tv.Links = links
			tri.Vertexes[i] = tv
		}
		d.smd.Triangles = append(d.smd.Triangles, tri)
	}
}

func (d *decoder) Scan() (string, error) {
	for {
		if !d.sc.Scan() {
			return "", io.EOF
		}
		line := d.sc.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		return strings.TrimSpace(line), nil
	}
}
