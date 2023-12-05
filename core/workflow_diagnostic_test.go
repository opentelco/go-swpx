package core

import (
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/analysispb"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"github.com/stretchr/testify/assert"
)

func Test_analyzeLink(t *testing.T) {

	t.Run("empty dataset", func(t *testing.T) {

		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{},
			},
			{
				NetworkElement: &networkelementpb.Element{},
			},
			{
				NetworkElement: &networkelementpb.Element{},
			},
		}

		_, err := analyzeLink("GigabitEthernet0/0/1", data)
		assert.Error(t, err)
	})

	t.Run("Link is up", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_up,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_up,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_up,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
		}

		report, err := analyzeLink("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_OK, report[0].Result)

	})

	t.Run("Link is down", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_up,
						},
					},
				},
			},
		}

		report, err := analyzeLink("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
	})

	t.Run("Link is shut", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
		}

		report, err := analyzeLink("GigabitEthernet0/0/1", data)
		assert.Contains(t, report[0].Note, "Link has been shut throughout")
		assert.NoError(t, err)
	})

	t.Run("Link is flapping", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_up,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_down,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description:       "GigabitEthernet0/0/1",
							OperationalStatus: networkelementpb.Port_up,
							AdminStatus:       networkelementpb.Port_down,
						},
					},
				},
			},
		}

		report, err := analyzeLink("GigabitEthernet0/0/1", data)
		assert.Contains(t, report[0].Note, "Link has been changing state under the")
		assert.NoError(t, err)
	})
}

func TestAnalyzeTransceiver(t *testing.T) {
	t.Run("empty dataset", func(t *testing.T) {

		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{},
			},
			{
				NetworkElement: &networkelementpb.Element{},
			},
			{
				NetworkElement: &networkelementpb.Element{},
			},
		}

		_, err := analyzeTransceiver("GigabitEthernet0/0/1", data)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no port found for GigabitEthernet0/0/1")
	})

	t.Run("rx and tx AVG is below threshold", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -30.32,
										Tx: -23.23,
									},
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -30.2,
										Tx: -19.23,
									},
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -40.00,
										Tx: -22.23,
									},
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeTransceiver("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[0].Result)
		assert.Contains(t, report[0].Note, "below threshold")
		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[1].Result)
		assert.Contains(t, report[1].Note, "below threshold")
	})
	t.Run("rx and tx AVG is within threshold", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -13.32,
										Tx: -13.23,
									},
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -13.1,
										Tx: -14.1,
									},
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Transceiver: &networkelementpb.Transceiver{
								Stats: []*networkelementpb.Transceiver_Statistics{
									{
										Rx: -15.00,
										Tx: -13.23,
									},
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeTransceiver("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_OK, report[0].Result)
		assert.Contains(t, report[0].Note, "within threshold")
		assert.Equal(t, analysispb.Analysis_RESULT_OK, report[1].Result)
		assert.Contains(t, report[1].Note, "within threshold")
	})
}

func TestAnalyzeErrors(t *testing.T) {
	t.Run("input errors increasing", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 1,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 3,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_ERROR, report[0].Result)
		assert.Contains(t, report[0].Note, "increased")
		assert.Contains(t, report[0].Note, "input")

	})

	t.Run("output errors increasing", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 1,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 3,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_ERROR, report[0].Result)
		assert.Contains(t, report[0].Note, "increased")
		assert.Contains(t, report[0].Note, "output")

	})

	t.Run("output and input errors increasing", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 1,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 1,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 3,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 2,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 8,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))

		assert.Equal(t, analysispb.Analysis_RESULT_ERROR, report[0].Result)
		assert.Contains(t, report[0].Note, "increased")
		assert.Contains(t, report[0].Note, "input")

		assert.Equal(t, analysispb.Analysis_RESULT_ERROR, report[1].Result)
		assert.Contains(t, report[1].Note, "increased")
		assert.Contains(t, report[1].Note, "output")

	})

	t.Run("output & input errors ", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[0].Result)
		assert.Contains(t, report[0].Note, "errors on the port")

		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[1].Result)
		assert.Contains(t, report[1].Note, "errors on the port")

	})

	t.Run("output  errors ", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 10,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[0].Result)
		assert.Contains(t, report[0].Note, "errors on the port")

	})

	t.Run("no errors ", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 0,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 0,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Errors: 0,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeErrors("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_OK, report[0].Result)
		assert.Contains(t, report[0].Note, "no errors on the port")

	})

}

func TestAnalyzeTraffic(t *testing.T) {
	t.Run("no traffic", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeTraffic("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_WARNING, report[0].Result)
		assert.Contains(t, report[0].Note, "not enough")

	})

	t.Run("traffic", func(t *testing.T) {
		data := []*corepb.PollResponse{
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{
							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 300,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 400,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 450,
								},
							},
						},
					},
				},
			},
			{
				NetworkElement: &networkelementpb.Element{
					Interfaces: []*networkelementpb.Port{
						{

							Description: "GigabitEthernet0/0/1",
							Stats: &networkelementpb.Port_Statistics{
								Output: &networkelementpb.Port_Statistics_Metrics{
									Packets: 600,
								},
								Input: &networkelementpb.Port_Statistics_Metrics{
									Packets: 700,
								},
							},
						},
					},
				},
			},
		}

		report, err := analyzeTraffic("GigabitEthernet0/0/1", data)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(report))
		assert.Equal(t, analysispb.Analysis_RESULT_OK, report[0].Result)
		assert.Contains(t, report[0].Note, "increased")
		assert.Contains(t, report[1].Note, "increased")

	})
}

func TestAverage(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"one element", []float64{5.0}, 5.0},
		{"multiple elements", []float64{1.0, 2.0, 3.0, 4.0}, 2.5},
		{"negative values", []float64{-1.0, -2.0, -3.0, -4.0}, -2.5},
		{"mixed values", []float64{-2.0, 1.0, 3.0}, 0.6666666666666666},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := average(tc.slice)
			if got != tc.expected {
				t.Errorf("average(%v) = %v; want %v", tc.slice, got, tc.expected)
			}
		})
	}
}

func TestFloat64ToString(t *testing.T) {
	slice := []float64{1.0, 2.0, 3.0}
	strSlice := float64ToString(slice)
	assert.Equal(t, []string{"1.000", "2.000", "3.000"}, strSlice)
}

func TestHasIncreasingErrors(t *testing.T) {
	t.Run("is increasing", func(t *testing.T) {
		data := []int64{231, 240, 260}
		assert.True(t, hasIncreasingErrors(data))
	})

	t.Run("is not increasing", func(t *testing.T) {
		data := []int64{231, 231, 231}
		assert.False(t, hasIncreasingErrors(data))
	})

	t.Run("is one", func(t *testing.T) {
		data := []int64{231}
		assert.False(t, hasIncreasingErrors(data))
	})

}
