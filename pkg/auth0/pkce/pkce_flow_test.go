package pkce

import (
	"net"
	"testing"
)

func TestFlow_BrowserURL(t *testing.T) {
	server := &localServer{
		listener: &fakeListener{
			addr: &net.TCPAddr{Port: 12345},
		},
	}

	type fields struct {
		server   *localServer
		clientID string
		state    string
	}
	type args struct {
		baseURL string
		params  BrowserParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				server: server,
				state:  "xy/z",
			},
			args: args{
				baseURL: "https://aserto-demo.us.auth0.com/authorize",
				params: BrowserParams{
					ClientID:      "Yurzr4ZMOIsc9YA3WRXLYurOyM1yqOU2",
					RedirectURI:   "http://localhost:3000",
					Scopes:        []string{"openid", "profile", "email"},
					CodeVerifier:  "0xCodeVerifier#",
					CodeChallenge: "0xCodeChallenge#",
				},
			},
			want:    "https://aserto-demo.us.auth0.com/authorize?audience=&client_id=Yurzr4ZMOIsc9YA3WRXLYurOyM1yqOU2&code_challenge=0xCodeChallenge%23&code_challenge_method=S256&prompt=login&redirect_uri=http%3A%2F%2Flocalhost%3A3000&response_type=code&scope=openid+profile+email&state=xy%2Fz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			flow := &Flow{
				server:   tc.fields.server,
				clientID: tc.fields.clientID,
				state:    tc.fields.state,
			}
			got, err := flow.BrowserURL(tc.args.baseURL, tc.args.params)
			if (err != nil) != tc.wantErr {
				t.Errorf("Flow.BrowserURL() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("Flow.BrowserURL() = %v, want %v", got, tc.want)
			}
		})
	}
}
