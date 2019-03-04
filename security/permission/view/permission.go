package view

type Permission string

const Configuration Permission = "Configuration"
const PartyCompanyConfiguration Permission = "PartyCompanyConfiguration"
const PartyClientConfiguration Permission = "PartyClientConfiguration"
const PartyUserConfiguration Permission = "PartyUserConfiguration"
const DeviceConfiguration Permission = "DeviceConfiguration"

const Dashboards Permission = "Dashboards"
const LiveTrackingDashboard Permission = "LiveTrackingDashboard"
const HistoricalTrackingDashboard Permission = "HistoricalTrackingDashboard"
