package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
)

type SampleHandler struct {
	SampleService service.ISampleService
}

func (sh *SampleHandler) GetdailyCollectionSamples(ctx context.Context, request *pb.SampleReceivedTime, response *pb.SampleCollectionList) error {
	startTime := request.Starttime
	endTime := request.Endtime

	sampleCollections, err := sh.SampleService.GetdailyCollectionSamples(startTime, endTime, ctx)
	if err != nil {
		return err
	}

	response.SampleCollection = sampleCollections

	return nil
}

func (sh *SampleHandler) GetdailyCheckNonReceivedSamples(ctx context.Context, request *pb.SampleReceivedTime, response *pb.Sample_NonReceivedList) error {
	startTime := request.Starttime
	endTime := request.Endtime

	sampleCollections, err := sh.SampleService.GetdailyCheckNonReceivedSamples(startTime, endTime, ctx)
	if err != nil {
		return err
	}

	response.Sample_NonReceived = sampleCollections

	return nil
}

func (sh *SampleHandler) GenerateSampleID(ctx context.Context, req *pb.EmptyRequest, resp *pb.GenerateSampleIdResponse) error {
	sampleIDResp, err := sh.SampleService.GenerateSampleID(&pb.EmptyRequest{}, ctx)
	if err != nil {
		return err
	}

	resp.SampleId = sampleIDResp.SampleId
	return nil
}

func (sh *SampleHandler) GenerateBarcodeForSampleID(ctx context.Context, req *pb.GenerateBarcodeForSampleIdRequest, resp *pb.GenerateBarcodeForSampleIdResponse) error {
	barcodeResp, err := sh.SampleService.GenerateBarcodeForSampleID(req, ctx)
	if err != nil {
		return err
	}

	resp.Barcode = barcodeResp.Barcode
	return nil
}

func (sh *SampleHandler) ListSampleTests(context.Context, *pb.SampleId, *pb.SampleTestId) error {
	return nil
}

func (sh *SampleHandler) ListSamples(context.Context, *pb.IdList, *pb.SampleList) error {
	return nil
}

func (sh *SampleHandler) GetSampleTests(context.Context, *pb.SampleIdList, *pb.SampleTestList) error {
	return nil
}

func (sh *SampleHandler) GetSampleTubes(context.Context, *pb.SampleId, *pb.Tubes) error {
	return nil
}

func (sh *SampleHandler) GetSampleTubesCount(context.Context, *pb.SampleId, *pb.Sample_Tubes_Counts_Response) error {
	return nil
}

func (sh *SampleHandler) GetMultiSampleTubesCount(context.Context, *pb.SampleIdList, *pb.GetMultiSampleTubesCountListResponse) error {
	return nil
}

func (sh *SampleHandler) ListSamplesAccessionID(context.Context, *pb.AccessionIdsList, *pb.SampleList) error {
	return nil
}

func (sh *SampleHandler) GetSampleReceiveRecords(context.Context, *pb.GetSampleReceiveRecordsRequest, *pb.GetSampleReceiveRecordsResponse) error {
	return nil
}

func (sh *SampleHandler) GetSampleReceiveRecordsBatch(context.Context, *pb.GetSampleReceiveRecordsRequestList, *pb.GetSampleReceiveRecordsResponseMap) error {
	return nil
}

func (sh *SampleHandler) ModifySampleReceiveRecord(context.Context, *pb.ModifySampleReceiveRecordRequest, *pb.ModifySampleReceiveRecordResponse) error {
	return nil
}

func (sh *SampleHandler) ReceiveSampleTubes(context.Context, *pb.ReceiveSampleTubesRequestStaging, *pb.ReceiveSampleTubesResponse) error {
	return nil
}

func (sh *SampleHandler) GetSampleTestsInstrument(context.Context, *pb.SampleIdList, *pb.SampleTestInstrumentList) error {
	return nil
}

func (sh *SampleHandler) GetTubeSampleTypeInfoViaTubeTypeSymbol(context.Context, *pb.GetSampleTypeTubeTypeRequest, *pb.GetTubeSampleTypeInfoViaTubeTypeSymbolResponseMessage) error {
	return nil
}

func (sh *SampleHandler) GetSampleTubeTypeInfoViaSampleTypeCode(context.Context, *pb.GetSampleTypeBySampleTypeCodeRequest, *pb.SampleTybeDetailsWithTubes) error {
	return nil
}

func (sh *SampleHandler) GetSampleTubeTypeInfoViaSampleEmunCode(context.Context, *pb.GetSampleTypeBySampleTypeEmunRequest, *pb.SampleTybeDetailsWithTubes) error {
	return nil
}

func (sh *SampleHandler) ListSampleMininumInfo(context.Context, *pb.AccessionIdsList, *pb.SampleListMininum) error {
	return nil
}

func (sh *SampleHandler) New_SearchSamples(context.Context, *pb.New_SearchSamplesRequest, *pb.New_SearchSamplesResponse) error {
	return nil
}

func (sh *SampleHandler) GetSampleEarilestCollectionAndReceiveTime(context.Context, *pb.GetSampleEarilestCollectionAndReceiveTimeRequest, *pb.GetSampleEarilestCollectionAndReceiveTimeResponse) error {
	return nil
}

func (sh *SampleHandler) CheckSampleAttributes(context.Context, *pb.CheckSamplesAttributesRequest, *pb.CheckSamplesAttributesResponse) error {
	return nil
}

func (sh *SampleHandler) BatchCheckSampleAttributes(context.Context, *pb.BatchCheckSamplesAttributesRequest, *pb.BatchCheckSamplesAttributesResponse) error {
	return nil
}

func (sh *SampleHandler) New_BatchCheckSampleAttributes(context.Context, *pb.New_BatchCheckSamplesAttributesRequest, *pb.New_BatchCheckSamplesAttributesResponse) error {
	return nil
}

func (sh *SampleHandler) FuzzySearchPhlebotomists(context.Context, *pb.FuzzySearchPhlebotomistsRequest, *pb.PhlebotomistsResponse) error {
	return nil
}
