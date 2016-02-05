Get-ChildItem -Filter *.nuspec -Recurse |
Foreach-Object{
    &$PsScriptRoot\NuGet.exe pack $_.FullName -Version 0.0.$env:BUILD_NUMBER
    #&$PsScriptRoot\NuGet.exe pack $_.FullName -Version 0.0.1.$env:BUILD_NUMBER
}
