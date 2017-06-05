package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/ant0ine/go-json-rest/rest"

	"common"
	"engine"
	"logger"
)

type statisticsServer struct {
	gateway *engine.Gateway
}

func newStatisticsServer(gateway *engine.Gateway) (*statisticsServer, error) {
	return &statisticsServer{
		gateway: gateway,
	}, nil
}

func (s *statisticsServer) Api() (*rest.Api, error) {
	router, err := rest.MakeRouter(
		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/plugins/#pluginName/indicators"),
			s.retrievePluginIndicatorNames),
		rest.Get(common.PrefixAPIVersion(
			"/pipelines/#pipelineName/plugins/#pluginName/indicators/#indicatorName/value"),
			s.retrievePluginIndicatorValue),
		rest.Get(common.PrefixAPIVersion(
			"/pipelines/#pipelineName/plugins/#pluginName/indicators/#indicatorName/desc"),
			s.retrievePluginIndicatorDesc),

		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/indicators"),
			s.retrievePipelineIndicatorNames),
		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/indicators/#indicatorName/value"),
			s.retrievePipelineIndicatorValue),
		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/indicators/#indicatorName/desc"),
			s.retrievePipelineIndicatorDesc),

		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/task/indicators"),
			s.retrievePipelineTaskIndicatorNames),
		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/task/indicators/#indicatorName/value"),
			s.retrievePipelineTaskIndicatorValue),
		rest.Get(common.PrefixAPIVersion("/pipelines/#pipelineName/task/indicators/#indicatorName/desc"),
			s.retrievePipelineTaskIndicatorDesc),

		rest.Get(common.PrefixAPIVersion("/gateway/uptime"), s.retrieveGatewayUpTime),
		rest.Get(common.PrefixAPIVersion("/gateway/rusage"), s.retrieveGatewaySysResUsage),
		rest.Get(common.PrefixAPIVersion("/gateway/loadavg"), s.retrieveGatewaySysAverageLoad),
	)

	if err != nil {
		logger.Errorf("[make router for staticstics server failed: %v]", err)
		return nil, err
	}

	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	api.SetApp(router)

	return api, nil
}

func (s *statisticsServer) retrievePluginIndicatorNames(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve plugin indicator names]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pluginName, err := url.QueryUnescape(r.PathParam("pluginName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PluginIndicatorNames(pluginName)
	// Returns with stable order
	sort.Strings(indicatorNames)

	w.WriteJson(&IndicatorNamesRetrieveResponse{
		Names: indicatorNames,
	})
	w.WriteHeader(http.StatusOK)

	logger.Debugf("[indicator names of plugin %s in pipeline %s returned]", pluginName, pipelineName)
}

func (s *statisticsServer) retrievePluginIndicatorValue(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve plugin indicator value]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pluginName, err := url.QueryUnescape(r.PathParam("pluginName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PluginIndicatorNames(pluginName)
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorValue, err := statistics.PluginIndicatorValue(pluginName, indicatorName)
	if err != nil {
		msg := fmt.Sprintf("evaluate indicator %s value failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorValueRetrieveResponse{
			Value: indicatorValue,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator value of plugin %s in pipeline %s returned]", pluginName, pipelineName)
	}
}

func (s *statisticsServer) retrievePluginIndicatorDesc(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve plugin indicator description]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pluginName, err := url.QueryUnescape(r.PathParam("pluginName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PluginIndicatorNames(pluginName)
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorDesc, err := statistics.PluginIndicatorDescription(pluginName, indicatorName)
	if err != nil {
		msg := fmt.Sprintf("describe indicator %s failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorDescriptionRetrieveResponse{
			Description: indicatorDesc,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator description of plugin %s in pipeline %s returned]", pluginName, pipelineName)
	}
}

func (s *statisticsServer) retrievePipelineIndicatorNames(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline indicator names]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PipelineIndicatorNames()
	// Returns with stable order
	sort.Strings(indicatorNames)

	w.WriteJson(&IndicatorNamesRetrieveResponse{
		Names: indicatorNames,
	})
	w.WriteHeader(http.StatusOK)

	logger.Debugf("[indicator names of pipeline %s returned]", pipelineName)
}

func (s *statisticsServer) retrievePipelineIndicatorValue(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline indicator value]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PipelineIndicatorNames()
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorValue, err := statistics.PipelineIndicatorValue(indicatorName)
	if err != nil {
		msg := fmt.Sprintf("evaluate indicator %s value failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorValueRetrieveResponse{
			Value: indicatorValue,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator value of pipeline %s returned]", pipelineName)
	}
}

func (s *statisticsServer) retrievePipelineIndicatorDesc(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline indicator description]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.PipelineIndicatorNames()
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorDesc, err := statistics.PipelineIndicatorDescription(indicatorName)
	if err != nil {
		msg := fmt.Sprintf("describe indicator %s failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorDescriptionRetrieveResponse{
			Description: indicatorDesc,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator description of pipeline %s returned]", pipelineName)
	}
}

func (s *statisticsServer) retrievePipelineTaskIndicatorNames(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline task indicator names]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.TaskIndicatorNames()
	// Returns with stable order
	sort.Strings(indicatorNames)

	w.WriteJson(&IndicatorNamesRetrieveResponse{
		Names: indicatorNames,
	})
	w.WriteHeader(http.StatusOK)

	logger.Debugf("[indicator names of task in pipeline %s returned]", pipelineName)
}

func (s *statisticsServer) retrievePipelineTaskIndicatorValue(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline task indicator value]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.TaskIndicatorNames()
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorValue, err := statistics.TaskIndicatorValue(indicatorName)
	if err != nil {
		msg := fmt.Sprintf("evaluate indicator %s value failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorValueRetrieveResponse{
			Value: indicatorValue,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator value of task in pipeline %s returned]", pipelineName)
	}
}

func (s *statisticsServer) retrievePipelineTaskIndicatorDesc(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve pipeline task indicator description]")

	pipelineName, err := url.QueryUnescape(r.PathParam("pipelineName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indicatorName, err := url.QueryUnescape(r.PathParam("indicatorName"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics := s.gateway.Model().StatRegistry().GetPipelineStatistics(pipelineName)
	if statistics == nil {
		msg := fmt.Sprintf("pipeline %s statistics not found", pipelineName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorNames := statistics.TaskIndicatorNames()
	if !common.StrInSlice(indicatorName, indicatorNames) {
		msg := fmt.Sprintf("indicator %s not found", indicatorName)
		rest.Error(w, msg, http.StatusNotFound)
		logger.Warnf("[%s]", msg)
		return
	}

	indicatorDesc, err := statistics.TaskIndicatorDescription(indicatorName)
	if err != nil {
		msg := fmt.Sprintf("describe indicator %s failed", indicatorName)
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(&IndicatorDescriptionRetrieveResponse{
			Description: indicatorDesc,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[indicator description of task in pipeline %s returned]", pipelineName)
	}
}

func (s *statisticsServer) retrieveGatewayUpTime(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve gateway uptime]")
	w.WriteJson(&GatewayUpTimeRetrieveResponse{
		UpTime: s.gateway.UpTime(),
	})
	w.WriteHeader(http.StatusOK)
	logger.Debugf("[gateway uptime returned]")
}

func (s *statisticsServer) retrieveGatewaySysResUsage(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve gateway system resource usage]")

	usage, err := s.gateway.SysResUsage()
	if err != nil {
		msg := fmt.Sprintf("get gateway system resource usage failed")
		rest.Error(w, msg, http.StatusInternalServerError)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(usage)
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[gateway system resource usage returned]")
	}
}

func (s *statisticsServer) retrieveGatewaySysAverageLoad(w rest.ResponseWriter, r *rest.Request) {
	logger.Debugf("[retrieve gateway system average load]")

	load1, load5, load15, err := s.gateway.SysAverageLoad()
	if err != nil {
		msg := fmt.Sprintf("get gateway system average load failed")
		rest.Error(w, msg, http.StatusForbidden)
		logger.Warnf("[%s: %v]", msg, err)
	} else {
		w.WriteJson(struct {
			Load1  float64 `json:"load1"`
			Load5  float64 `json:"load5"`
			Load15 float64 `json:"load15"`
		}{
			Load1:  load1,
			Load5:  load5,
			Load15: load15,
		})
		w.WriteHeader(http.StatusOK)
		logger.Debugf("[gateway system average load returned]")
	}
}
