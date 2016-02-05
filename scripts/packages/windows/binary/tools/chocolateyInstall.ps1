function New-NetFirewallRule
{
	param
	(
		[String]$DisplayName,
		[String]$Direction = "Inbound",
		[String]$LocalPort,
		[String]$Protocol  = "TCP",
		[String]$Action    = "Allow"
	)

	$dir = switch ($Direction) 
	{ 
		"Inbound"  {"in"} 
		"Outbound" {"out"}
	}

	$argumentList = @(
		'advfirewall', 'firewall', 'add', 'rule'
		"name=$DisplayName"
		"dir=$dir"
		"protocol=$Protocol"
		"localport=$LocalPort" #${env:Logstash.Port}"
		"action=$Action"
	)
  
    Write-Host "Adding firewall rule: " $DisplayName " Port:" $LocalPort " Direction:" $Direction " Protocol:" $Protocol " Action:" $Action
	$process = Start-Process netsh -ArgumentList $argumentList -Wait -NoNewWindow -PassThru
	if ($process.ExitCode -ne 0) { throw "Error adding new firewall rule"}
}

function Update-Logstash-Service-Parameters
{ 
	param
	(
		[String]$Servicename
	)

    Get-Service $Servicename
    if ( $? -eq "True" )
    {
        &$nssmexe status $Servicename 
        if ( $? -eq "False" )
        {
            &$nssmexe  stop $Servicename 
            #&$nssmexe  status  $Servicename 
        }

		Write-Host "Setting service path to " $appDirectory
		&$nssmexe set $servicename AppDirectory $appDirectory
    
		Write-Host "Setting service parameters to " $appParameters
		&$nssmexe set $servicename AppParameters $appParameters
    
		Write-Host "Setting service display name to " $displayName
		&$nssmexe set $servicename DisplayName $displayName
		
		Write-Host "Setting service description to " $displayName
		&$nssmexe set $servicename Description $displayName

        &$nssmexe  status  $servicename 
        if ( $? -eq "False" )
        {
            &$nssmexe  start $servicename 
            &$nssmexe  status  $servicename 
        }
    }
}


try {  
    $packageName        = "SEEK-Logstash-Aggregator"
	$logstashVersion	= "1.4.1"
    $logstashRemotePort = "9299"
    $esRemotePort       = "9200"
    $servicename        = "Logstash"
    $displayName        = "$packageName" + "-" + $logstashVersion
    $appDirectory       = "c:\logstash\bin"
    $appCmd             = $appDirectory + "\logstash.bat"

    # path has to be in unix format eg: c:/logstash/bin/conf/
    $configDirectory    = ( Resolve-Path("$PsScriptRoot\..\conf\") ) -replace '\\','/' 
    $appParameters      = "agent -f " + $configDirectory
    $nssmexe            = "nssm.exe" 

	New-NetFirewallRule -DisplayName $packageName -LocalPort $logstashRemotePort -Direction Inbound
	New-NetFirewallRule -DisplayName $packageName -LocalPort $esRemotePort -Direction Outbound

	Update-Logstash-Service-Parameters $servicename
} 
catch {
    throw $_.Exception
}
