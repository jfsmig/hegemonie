// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"encoding/json"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"os"
)

func (cli *ClientCLI) DoTemplateList(ctx context.Context, reg, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewTemplatesClient(cnx).ListTemplates(ctx,
			&proto.PaginatedStrQuery{Region: reg, Marker: marker})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			t := _cityTemplate{}
			x, e := out.Recv()
			if x != nil {
				t.importFrom(x)
			}
			return &t, e
		})
	})
}

func (cli *ClientCLI) DoTemplateCreate(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		cliTpl := emptyCityTemplate()
		if err := json.NewDecoder(os.Stdin).Decode(&cliTpl); err != nil {
			return errors.Trace(err)
		}
		tpl := cliTpl.exportTo()
		out, err := proto.NewTemplatesClient(cnx).CreateTemplate(ctx,
			&proto.CityTemplateReq{Region: reg, Tpl: tpl})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(out)
	})
}

func (cli *ClientCLI) DoTemplateUpdate(ctx context.Context, reg, tplID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		cliTpl := emptyCityTemplate()
		if err := json.NewDecoder(os.Stdin).Decode(&cliTpl); err != nil {
			return errors.Trace(err)
		}
		tpl := cliTpl.exportTo()
		tpl.Name = tplID
		out, err := proto.NewTemplatesClient(cnx).UpdateTemplate(ctx,
			&proto.CityTemplateReq{Region: reg, Id: tplID, Tpl: tpl})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(out)
	})
}

func (cli *ClientCLI) DoTemplateDelete(ctx context.Context, reg, tpl string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewTemplatesClient(cnx).DeleteTemplate(ctx,
			&proto.TemplateId{Region: reg, Id: tpl})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(out)
	})
}

// DoTemplateEmpty generates to os.Stdout a JSON representation of a valid but zeroed CityTemplate
func DoTemplateEmpty() error {
	return utils.DumpJSON(emptyCityTemplate())
}

// DoTemplateSet generates to os.Stdout a JSON representation of the CityTemplate loaded from os.Stdin
// with some fields altered.
// @param kv must be a sequence of string pairs, with the second value being an encoded integer
func DoTemplateSet(kv []string) error {
	if (len(kv) % 2) != 0 {
		return errors.BadRequestf("even number of parameter, not mappable to key/value")
	}
	tpl := emptyCityTemplate()
	if err := json.NewDecoder(os.Stdin).Decode(&tpl); err != nil {
		return errors.Trace(err)
	}

	for i := 0; i < len(kv); i += 2 {
		if err := tpl.set(kv[i], kv[i+1]); err != nil {
			return errors.Trace(err)
		}
	}
	return utils.DumpJSON(tpl)
}

// DoTemplateDel
func DoTemplateDel(kv []string) error {
	if (len(kv) % 2) != 0 {
		return errors.BadRequestf("even number of parameter, not mappable to key/value")
	}
	tpl := emptyCityTemplate()
	if err := json.NewDecoder(os.Stdin).Decode(&tpl); err != nil {
		return errors.Trace(err)
	}

	for i := 0; i < len(kv); i += 2 {
		if err := tpl.del(kv[i], kv[i+1]); err != nil {
			return errors.Trace(err)
		}
	}
	return utils.DumpJSON(tpl)
}

// DoTemplateAdd
func DoTemplateAdd(kv []string) error {
	if (len(kv) % 2) != 0 {
		return errors.BadRequestf("even number of parameter, not mappable to key/value")
	}
	tpl := emptyCityTemplate()
	if err := json.NewDecoder(os.Stdin).Decode(&tpl); err != nil {
		return errors.Trace(err)
	}

	for i := 0; i < len(kv); i += 2 {
		if err := tpl.add(kv[i], kv[i+1]); err != nil {
			return errors.Trace(err)
		}
	}
	return utils.DumpJSON(tpl)
}
