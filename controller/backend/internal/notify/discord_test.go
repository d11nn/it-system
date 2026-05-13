package notify

import "testing"

func TestFormatNfPrSummary(t *testing.T) {
	testCases := []struct {
		name     string
		nfPrList []NfPrResult
		expected string
	}{
		{
			name:     "empty",
			nfPrList: nil,
			expected: "No PRs",
		},
		{
			name: "single PR",
			nfPrList: []NfPrResult{
				{NfName: "amf", PR: 213},
			},
			expected: "amf #213",
		},
		{
			name: "multiple PRs under limit",
			nfPrList: []NfPrResult{
				{NfName: "amf", PR: 213},
				{NfName: "smf", PR: 200},
			},
			expected: "amf #213, smf #200",
		},
		{
			name: "multiple PRs over limit",
			nfPrList: []NfPrResult{
				{NfName: "access-and-mobility", PR: 213},
				{NfName: "session-management", PR: 200},
				{NfName: "network-repository", PR: 33},
			},
			expected: "access-and-mobility #213 +2 more",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatNfPrSummary(tc.nfPrList)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestFormatNfPrDetails(t *testing.T) {
	actual := formatNfPrDetails([]NfPrResult{
		{NfName: "amf", PR: 213},
		{NfName: "smf", PR: 200},
	})
	expected := "- amf #213\n- smf #200"

	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}
