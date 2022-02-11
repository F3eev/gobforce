package nmapxml


import "encoding/xml"

// Run contains the main data of an nmap scan
type Run struct {
	XMLName          xml.Name `xml:"nmaprun" json:"-"`
	Scanner          string   `xml:"scanner,attr" json:"scanner"`
	Args             string   `xml:"args,attr" json:"args"`
	Start            string   `xml:"start,attr" json:"start"`
	StartStr         string   `xml:"startstr,attr" json:"startstr"`
	Version          string   `xml:"version,attr" json:"version"`
	XMLOutputVersion string   `xml:"xmloutputversion,attr" json:"xmloutputversion"`
	ScanInfo         ScanInfo
	Verbose          Verbose
	Debugging        Debugging
	Host             []Host `xml:"host" json:"host"`
}

// ScanInfo contains metadata of an nmap scan
type ScanInfo struct {
	XMLName     xml.Name `xml:"scaninfo" json:"-"`
	Type        string   `xml:"type,attr" json:"type"`
	Protocol    string   `xml:"protocol,attr" json:"protocol"`
	Numservices string   `xml:"numservices,attr" json:"numservices"`
	Services    string   `xml:"services,attr" json:"services"`
}

// Verbose contains verbosity level of an nmap scan
type Verbose struct {
	XMLName xml.Name `xml:"verbose" json:"-"`
	Level   string   `xml:"level,attr" json:"level"`
}

// Debugging contains debugging level of an nmap scan
type Debugging struct {
	XMLName xml.Name `xml:"debugging" json:"-"`
	Level   string   `xml:"level,attr" json:"level"`
}

// Host contains a scanned host
type Host struct {
	XMLName   xml.Name `xml:"host" json:"-"`
	StartTime string   `xml:"starttime,attr" json:"starttime"`
	EndTime   string   `xml:"endtime,attr" json:"endtime"`
	Status    Status
	Address   Address
	Hostnames Hostnames
	Ports     Ports
	Times     Times
}

// Status contains status of a scanned host
type Status struct {
	XMLName   xml.Name `xml:"status" json:"-"`
	State     string   `xml:"state,attr" json:"state"`
	Reason    string   `xml:"reason,attr" json:"reason"`
	ReasonTTL string   `xml:"reason_ttl,attr" json:"reason_ttl"`
}

// Address contains IP address of a scanned host
type Address struct {
	XMLName  xml.Name `xml:"address" json:"-"`
	Addr     string   `xml:"addr,attr" json:"addr"`
	AddrType string   `xml:"addrtype,attr" json:"addrtype"`
}

// Hostnames wraps an array of Hostname nodes
type Hostnames struct {
	XMLName  xml.Name `xml:"hostnames" json:"-"`
	Hostname Hostname
}

// Hostname contains hostname of a scanned host
type Hostname struct {
	XMLName xml.Name `xml:"hostname" json:"-"`
	Name    *string  `xml:"name,attr" json:"name,omitempty"`
	Type    *string  `xml:"type,attr" json:"type,omitempty"`
}

// Ports contains all ports scanned of a scanned host
type Ports struct {
	XMLName    xml.Name   `xml:"ports" json:"-"`
	ExtraPorts ExtraPorts `xml:"extraports" json:"extraports"`
	Port       *[]Port    `xml:"port" json:"port,omitempty"`
}

// ExtraPorts contains any non-open ports of a scanned host
type ExtraPorts struct {
	XMLName      xml.Name      `xml:"extraports" json:"-"`
	State        string        `xml:"state,attr" json:"state"`
	Count        string        `xml:"count,attr" json:"count"`
	ExtraReasons []ExtraReason `xml:"extrareasons" json:"extrareasons"`
}

// ExtraReason contains metadata of extra ports
type ExtraReason struct {
	XMLName xml.Name `xml:"extrareasons" json:"-"`
	Reason  string   `xml:"reason,attr" json:"reasons"`
	Count   string   `xml:"count,attr" json:"count"`
}

// Port contains an open port of a scanned host
type Port struct {
	XMLName  xml.Name `xml:"port" json:"-"`
	Protocol string   `xml:"protocol,attr" json:"protocol"`
	PortID   string   `xml:"portid,attr" json:"portid"`
	State    State    `xml:"state" json:"state"`
	Service  Service  `xml:"service" json:"service"`
	Script   *Script  `xml:"script" json:"script,omitempty"`
}

// State contains state information of a port
type State struct {
	XMLName   xml.Name `xml:"state" json:"-"`
	State     string   `xml:"state,attr" json:"state"`
	Reason    string   `xml:"reason,attr" json:"reason"`
	ReasonTTL string   `xml:"reason_ttl,attr" json:"reason_ttl"`
}

// Service contains service information of a port
type Service struct {
	XMLName    xml.Name `xml:"service" json:"-"`
	Name       string   `xml:"name,attr" json:"name"`
	Product    *string  `xml:"product,attr" json:"product,omitempty"`
	DeviceType *string  `xml:"devicetype,attr" json:"devicetype,omitempty"`
	ServiceFP  *string  `xml:"servicefp,attr" json:"servicefp,omitempty"`
	Tunnel     *string  `xml:"tunnel,attr" json:"tunnel,omitempty"`
	Method     string   `xml:"method,attr" json:"method"`
	Conf       string   `xml:"conf,attr" json:"conf"`
	CPE        *string  `xml:"cpe" json:"cpe,omitempty"`
}

// Script contains metadata of script ran on a port
type Script struct {
	XMLName    xml.Name `xml:"script" json:"-"`
	ID         string   `xml:"id,attr" json:"id"`
	Output     string   `xml:"output,attr" json:"output"`
	Table      *[]Table `xml:"table" json:"table,omitempty"`
	ScriptData *[]Elem  `xml:"elem"`
}

// Elem contains output of a script
type Elem struct {
	Key   string `xml:"key,attr" json:"key,omitempty"`
	Value string `xml:",innerxml" json:"value,omitempty"`
}

// Table contains a table of elems
type Table struct {
	XMLName   xml.Name `xml:"table" json:"-"`
	Key       *string  `xml:"key,attr" json:"key,omitempty"`
	Table     *Table   `json:",omitempty"`
	TableData *[]Elem  `xml:"elem" json:",omitempty"`
}

// Times contains latency information of a scanned host
//
// Srtt == Smoothed Averaged Round Trip Time
//
// RTT == Round Trip Time
//
// To == ?
type Times struct {
	XMLName xml.Name `xml:"times" json:"-"`
	Srtt    string   `xml:"srtt,attr" json:"srtt"`
	RttVar  string   `xml:"rttvar,attr" json:"rttvar"`
	To      string   `xml:"to,attr" json:"to"`
}

// Runstats contains final information of a scan
type Runstats struct {
	XMLName  xml.Name `xml:"runstats" json:"-"`
	Finished Finished `xml:"finished" json:"finished"`
	Hosts    Hosts    `xml:"hosts" json:"hosts"`
}

// Finished contains total time a scan took
type Finished struct {
	XMLName xml.Name `xml:"finished" json:"-"`
	Time    string   `xml:"time,attr" json:"time"`
	Timestr string   `xml:"timestr,attr" json:"timestr"`
	Elapsed string   `xml:"elapsed,attr" json:"elapsed"`
	Summary string   `xml:"summary,attr" json:"summary"`
	Exit    string   `xml:"exit,attr" json:"exit"`
}

// Hosts contains counts of hosts up vs down and a count of total hosts of a scan
type Hosts struct {
	XMLName xml.Name `xml:"hosts" json:"-"`
	Up      string   `xml:"up,attr" json:"up"`
	Down    string   `xml:"down,attr" json:"down"`
	Total   string   `xml:"total,attr" json:"total"`
}
