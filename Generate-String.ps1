$KeyStr   = ""
$UrlToZip = ""

$Aes = [System.Security.Cryptography.Aes]::Create()
$Aes.Mode = [System.Security.Cryptography.CipherMode]::CBC
$Aes.Padding = [System.Security.Cryptography.PaddingMode]::PKCS7
$Aes.Key = [System.Text.Encoding]::UTF8.GetBytes($KeyStr)
$Aes.GenerateIV()
$IV = $Aes.IV

$Encryptor = $Aes.CreateEncryptor()
$PlainBytes = [System.Text.Encoding]::UTF8.GetBytes($UrlToZip)
$EncryptedBytes = $Encryptor.TransformFinalBlock($PlainBytes, 0, $PlainBytes.Length)

$Result = $IV + $EncryptedBytes

$Base64 = [Convert]::ToBase64String($Result)
$StringSafe = $Base64.Replace('+', '-').Replace('/', '_').TrimEnd('=')

Write-Host "Generated String:" -ForegroundColor Green
Write-Host "$StringSafe"
