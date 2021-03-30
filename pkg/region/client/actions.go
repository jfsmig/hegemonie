// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"encoding/json"
	region "github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

// ClientCLI gathers the actions destined to be exposed at the CLI, to manage a region service.
type ClientCLI struct{}

func (cli *ClientCLI) DoCityBuild(ctx context.Context, regID, cityStrID, bStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, bID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		bID, err = strconv.ParseUint(bStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.BuildReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			BuildingType: bID,
		}
		_, err = proto.NewCityClient(cnx).Build(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityStudy(ctx context.Context, regID, cityStrID, sStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, skillID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		skillID, err = strconv.ParseUint(sStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.StudyReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			KnowledgeType: skillID,
		}
		_, err = proto.NewCityClient(cnx).Study(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityTrain(ctx context.Context, regID, cityStrID, uStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, unitID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		unitID, err = strconv.ParseUint(uStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.TrainReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			UnitType: unitID,
		}
		_, err = proto.NewCityClient(cnx).Train(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityShow(ctx context.Context, regID, cityStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID uint64
		var cityView *proto.CityView

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		id := proto.CityId{
			Region: regID, City: cityID,
		}

		cityView, err = proto.NewCityClient(cnx).ShowAll(ctx, &id)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(cityView)
	})
}

// DoCreateRegion triggers the synchronous creation of a region with the given name, modeled on the named map.
func (cli *ClientCLI) DoCreateRegion(ctx context.Context, regID, mapID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).CreateRegion(ctx, &proto.RegionCreateReq{MapName: mapID, Name: regID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, regID, "created")
	})
}

// DoListRegions dumps to os.Stdout a JSON stream of the known regions, sorted by name
func (cli *ClientCLI) DoListRegions(ctx context.Context, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewAdminClient(cnx).ListRegions(ctx, &proto.RegionListReq{NameMarker: marker})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// DoRegionMovement triggers one round of armies movement on all the cities of the named region
func (cli *ClientCLI) DoRegionMovement(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Move(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Moved")
	})
}

// DoRegionProduction triggers the production of resources on all the cities of the named region
func (cli *ClientCLI) DoRegionProduction(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Produce(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Produced")
	})
}

// DoRegionPushStats triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionPushStats(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).PushStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Done")
	})
}

type _resourcesAbs struct {
	R0 uint64 `json:"r0"`
	R1 uint64 `json:"r1"`
	R2 uint64 `json:"r2"`
	R3 uint64 `json:"r3"`
	R4 uint64 `json:"r4"`
	R5 uint64 `json:"r5"`
}

type _cityStats struct {
	// Identifier
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// Gauges
	StockCapacity  _resourcesAbs `json:"stockCapacity"`
	StockUsage     _resourcesAbs `json:"stockUsage"`
	ScoreBuilding  uint64        `json:"scoreBuilding"`
	ScoreKnowledge uint64        `json:"scoreKnowledge"`
	ScoreArmy      uint64        `json:"scoreArmy"`
	// Counters
	ResourceProduced _resourcesAbs `json:"resourceProduced"`
	ResourceSent     _resourcesAbs `json:"resourceSent"`
	ResourceReceived _resourcesAbs `json:"resourceReceived"`
	TaxSent          _resourcesAbs `json:"taxSent"`
	TaxReceived      _resourcesAbs `json:"taxReceived"`
	Moves            uint64        `json:"moves"`
	UnitRaised       uint64        `json:"unitRaised"`
	UnitLost         uint64        `json:"unitLost"`
	FightJoined      uint64        `json:"fightJoined"`
	FightLeft        uint64        `json:"fightLeft"`
	FightWon         uint64        `json:"fightWon"`
	FightLost        uint64        `json:"fightLost"`
}

func resProto2Json(in *proto.ResourcesAbs) _resourcesAbs {
	return _resourcesAbs{
		R0: in.R0,
		R1: in.R1,
		R2: in.R2,
		R3: in.R3,
		R4: in.R4,
		R5: in.R5,
	}
}
func statsProto2Json(in *proto.CityStats) _cityStats {
	return _cityStats{
		ID:   in.Id,
		Name: in.Name,

		StockCapacity:  resProto2Json(in.StockCapacity),
		StockUsage:     resProto2Json(in.StockUsage),
		ScoreBuilding:  in.ScoreBuilding,
		ScoreKnowledge: in.ScoreKnowledge,
		ScoreArmy:      in.ScoreArmy,

		ResourceProduced: resProto2Json(in.ResourceProduced),
		ResourceSent:     resProto2Json(in.ResourceSent),
		ResourceReceived: resProto2Json(in.ResourceReceived),
		TaxSent:          resProto2Json(in.TaxSent),
		TaxReceived:      resProto2Json(in.TaxReceived),

		Moves:       in.Moves,
		UnitRaised:  in.UnitRaised,
		UnitLost:    in.UnitLost,
		FightJoined: in.FightJoined,
		FightLeft:   in.FightLeft,
		FightWon:    in.FightWon,
		FightLost:   in.FightLost,
	}
}

// DoRegionGetStats triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionGetStats(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewAdminClient(cnx).GetStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			itf, err := rep.Recv()
			if err != nil {
				return nil, err
			}
			// TODO(jfs): Map the protobuf-generated struct to a fac-simile
			// 			  with "omitempty" flags missing. There are probably
			// 			  better ways to achieve this.
			return statsProto2Json(itf), nil
		})
	})
}

// _publicCity is a variant of proto.PublicCity that doesn't omit empty
// fields. So that the value will me printed on os.Stdout
type _publicCity struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Alignment int32  `json:"alignment"`
	Chaos     int32  `json:"chaos"`
	Politics  uint32 `json:"politics"`
	Cult      uint32 `json:"cult"`
	Ethny     uint32 `json:"ethny"`
	Score     int64  `json:"score"`
}

type _cityTemplate struct {
	Public        _publicCity                `json:"public"`
	Stock         [region.ResourceMax]uint64 `json:"stock"`
	StockCapacity [region.ResourceMax]uint64 `json:"capacity"`
	Production    [region.ResourceMax]uint64 `json:"production"`
	BuildingTypes []uint64                   `json:"buildings"`
	SkillTypes    []uint64                   `json:"skills"`
	UnitTypes     []uint64                   `json:"units"`
}

func (t _cityTemplate) exportTo() *proto.CityTemplate {
	s := func(i int) uint64 { return t.Stock[i] }
	c := func(i int) uint64 { return t.StockCapacity[i] }
	p := func(i int) uint64 { return t.Production[i] }

	pub := t.Public.exportTo()

	tpl := proto.CityTemplate{Public: pub}
	tpl.Stock = &proto.ResourcesAbs{R0: s(0), R1: s(1), R2: s(2), R3: s(3), R4: s(4), R5: s(5)}
	tpl.StockCapacity = &proto.ResourcesAbs{R0: c(0), R1: c(1), R2: c(2), R3: c(3), R4: c(4), R5: c(5)}
	tpl.Production = &proto.ResourcesAbs{R0: p(0), R1: p(1), R2: p(2), R3: p(3), R4: p(4), R5: p(5)}
	tpl.BuildingTypes = t.BuildingTypes
	tpl.SkillTypes = t.SkillTypes
	tpl.UnitTypes = t.UnitTypes
	return &tpl
}

func emptyCityTemplate() _cityTemplate {
	return _cityTemplate{
		SkillTypes:    make([]uint64, 0),
		BuildingTypes: make([]uint64, 0),
		UnitTypes:     make([]uint64, 0),
	}
}

func (t _publicCity) exportTo() *proto.PublicCity {
	return &proto.PublicCity{
		Name:      t.Name,
		Ethny:     t.Ethny,
		Cult:      t.Cult,
		Politics:  t.Politics,
		Chaos:     t.Chaos,
		Alignment: t.Alignment,
	}
}

func (t *_publicCity) importFrom(tpl *proto.PublicCity) {
	t.Name = tpl.Name
	t.Ethny = tpl.Ethny
	t.Cult = tpl.Cult
	t.Chaos = tpl.Chaos
	t.Politics = tpl.Politics
	t.Alignment = tpl.Alignment
}

func (t *_cityTemplate) importFrom(tpl *proto.CityTemplate) {
	t.Public.importFrom(tpl.Public)

	t.Stock[0] = tpl.Stock.R0
	t.Stock[1] = tpl.Stock.R1
	t.Stock[2] = tpl.Stock.R2
	t.Stock[3] = tpl.Stock.R3
	t.Stock[4] = tpl.Stock.R4
	t.Stock[5] = tpl.Stock.R5

	t.StockCapacity[0] = tpl.StockCapacity.R0
	t.StockCapacity[1] = tpl.StockCapacity.R1
	t.StockCapacity[2] = tpl.StockCapacity.R2
	t.StockCapacity[3] = tpl.StockCapacity.R3
	t.StockCapacity[4] = tpl.StockCapacity.R4
	t.StockCapacity[5] = tpl.StockCapacity.R5

	t.Production[0] = tpl.Production.R0
	t.Production[1] = tpl.Production.R1
	t.Production[2] = tpl.Production.R2
	t.Production[3] = tpl.Production.R3
	t.Production[4] = tpl.Production.R4
	t.Production[5] = tpl.Production.R5

	empty := make([]uint64, 0)
	tpl.BuildingTypes = empty
	if tpl.BuildingTypes != nil {
		t.BuildingTypes = tpl.BuildingTypes
	}
	tpl.SkillTypes = empty
	if tpl.SkillTypes != nil {
		t.SkillTypes = tpl.SkillTypes
	}
	tpl.UnitTypes = empty
	if tpl.UnitTypes != nil {
		t.UnitTypes = tpl.UnitTypes
	}
}

func (t *_cityTemplate) set(k, vs string) error {
	i64, err := strconv.ParseInt(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "invalid value")
	}
	switch k {
	case "align":
		t.Public.Alignment = int32(i64)
	case "chaos":
		t.Public.Chaos = int32(i64)
	case "politics":
		t.Public.Politics = uint32(i64)
	case "cult":
		t.Public.Cult = uint32(i64)
	case "ethny":
		t.Public.Ethny = uint32(i64)

	case "stock.0", "s0":
		t.Stock[0] = uint64(i64)
	case "stock.1", "s1":
		t.Stock[1] = uint64(i64)
	case "stock.2", "s2":
		t.Stock[2] = uint64(i64)
	case "stock.3", "s3":
		t.Stock[3] = uint64(i64)
	case "stock.4", "s4":
		t.Stock[4] = uint64(i64)
	case "stock.5", "s5":
		t.Stock[5] = uint64(i64)

	case "capa.0", "c0":
		t.StockCapacity[0] = uint64(i64)
	case "capa.1", "c1":
		t.StockCapacity[1] = uint64(i64)
	case "capa.2", "c2":
		t.StockCapacity[2] = uint64(i64)
	case "capa.3", "c3":
		t.StockCapacity[3] = uint64(i64)
	case "capa.4", "c4":
		t.StockCapacity[4] = uint64(i64)
	case "capa.5", "c5":
		t.StockCapacity[5] = uint64(i64)

	case "prod.0", "p0":
		t.Production[0] = uint64(i64)
	case "prod.1", "p1":
		t.Production[1] = uint64(i64)
	case "prod.2", "p2":
		t.Production[2] = uint64(i64)
	case "prod.3", "p3":
		t.Production[3] = uint64(i64)
	case "prod.4", "p4":
		t.Production[4] = uint64(i64)
	case "prod.5", "p5":
		t.Production[5] = uint64(i64)
	}
	return nil
}

func (t *_cityTemplate) add(k, vs string) error {
	u64, err := strconv.ParseUint(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "Invalid value")
	}
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = append(t.BuildingTypes, u64)
	case "s", "skill", "skills":
		t.SkillTypes = append(t.SkillTypes, u64)
	case "u", "unit", "units":
		t.UnitTypes = append(t.UnitTypes, u64)
	}
	return nil
}

func (t *_cityTemplate) del(k, vs string) error {
	filter := func(tab []uint64, bad uint64) []uint64 {
		out := make([]uint64, 0)
		for _, v := range tab {
			if v != bad {
				out = append(out, v)
			}
		}
		return out
	}
	u64, err := strconv.ParseUint(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "Invalid value")
	}
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = filter(t.BuildingTypes, u64)
	case "s", "skill", "skills":
		t.SkillTypes = filter(t.SkillTypes, u64)
	case "u", "unit", "units":
		t.UnitTypes = filter(t.UnitTypes, u64)
	}
	return nil
}

// DoRegionGetScore triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionGetScores(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewAdminClient(cnx).GetScores(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			pc0, e := out.Recv()
			if e != nil {
				return nil, e
			}
			pc := _publicCity{}
			pc.importFrom(pc0)
			return &pc, nil
		})
	})
}

func (cli *ClientCLI) DoRegionGetBuildings(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListBuildings(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}

func (cli *ClientCLI) DoRegionGetSkills(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListKnowledges(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}

func (cli *ClientCLI) DoRegionGetUnits(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListUnits(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}

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
		tpl.Public.Name = tplID
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

func (cli *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Region()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
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
