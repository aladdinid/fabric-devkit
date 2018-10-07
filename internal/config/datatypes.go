/*
Copyright 2018 Aladdin Blockchain Technologies Ltd
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

// OrgSpec represents specification of Fabric type
type OrgSpec struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Anchor string `json:"anchor"`
}

// NetworkSpec represents specification of a Fabric network
type NetworkSpec struct {
	ScriptPath          string `json:"scriptpath"`
	ChaincodePath       string `json:"chaincodepath"`
	NetworkPath         string `json:"networkpath"`
	CryptoPath          string `json:"cryptopath"`
	ChannelArtefactPath string `json:"configtxpath"`
	Domain              string `json:"domain"`
	ChannelName         string `json:"channelname"`
	Consortium          string `json:"consortium"`
	OrganisationSpecs   []OrgSpec
}

// NewNetworkSpec instantiate a reference to a new spec
func NewNetworkSpec() *NetworkSpec {
	spec := new(NetworkSpec)

	spec.ChaincodePath = ChaincodePath()
	spec.ScriptPath = ScriptPath()
	spec.NetworkPath = NetworkPath()
	spec.CryptoPath = CryptoPath()
	spec.ChannelArtefactPath = ChannelArtefactPath()
	spec.Domain = Domain()
	spec.ChannelName = ChannelName()
	spec.Consortium = Consortium()
	spec.OrganisationSpecs = OrganizationSpecs()

	return spec
}
