package dpfm_api_input_reader

import (
	"data-platform-api-district-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToDistrict() *requests.District {
	data := sdc.District
	return &requests.District{
		District: data.District,
		Country:  data.Country,
	}
}
