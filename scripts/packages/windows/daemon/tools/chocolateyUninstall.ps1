function Remove-NetFirewallRule
{
	param
	(
		[String]$DisplayName
	)

	$argumentList = @(
		'advfirewall', 'firewall', 'del', 'rule'
		"name=$DisplayName"
	)
  
	Start-Process netsh -ArgumentList $argumentList -Wait -NoNewWindow
}

function DeRegister-Logstash-Service
{ 
	param
	(
		[String]$Servicename
	)

    Get-Service $Servicename
    if ( $? -eq "True" )
    {
        &$nssmexe status  $Servicename 
        if ( $? -eq "True" )
        {
            &$nssmexe stop $Servicename 
            #&$nssmexe remove $Servicename confirm
        }
        
        choco uninstall $Servicename 
    }
}

try {
    $packageName = "SEEK-Logstash-Aggregator"
    $servicename = "Logstash"
    $nssmexe     = "nssm.exe"

    Remove-NetFirewallRule -DisplayName $packageName
    DeRegister-Logstash-Service $servicename
}
catch {
    throw $_.Exception
}
