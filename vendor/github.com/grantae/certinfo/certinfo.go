package certinfo

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"
)

// Extra ASN1 OIDs that we may need to handle
var (
	oidEmailAddress                 = []int{1, 2, 840, 113549, 1, 9, 1}
	oidExtensionAuthorityInfoAccess = []int{1, 3, 6, 1, 5, 5, 7, 1, 1}
	oidNSComment                    = []int{2, 16, 840, 1, 113730, 1, 13}
)

// validity allows unmarshaling the certificate validity date range
type validity struct {
	NotBefore, NotAfter time.Time
}

// publicKeyInfo allows unmarshaling the public key
type publicKeyInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

// tbsCertificate allows unmarshaling of the "To-Be-Signed" principle portion
// of the certificate
type tbsCertificate struct {
	Version            int `asn1:"optional,explicit,default:1,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueID           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueID    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
}

// certUniqueIDs extracts the subject and issuer unique IDs which are
// byte strings. These are not common but may be present in x509v2 certificates
// or later under tags 1 and 2 (before x509v3 extensions).
func certUniqueIDs(tbsAsnData []byte) (issuerUniqueID, subjectUniqueID []byte, err error) {
	var tbs tbsCertificate
	rest, err := asn1.Unmarshal(tbsAsnData, &tbs)
	if err != nil {
		return nil, nil, err
	}
	if len(rest) > 0 {
		return nil, nil, asn1.SyntaxError{Msg: "trailing data"}
	}
	iuid := tbs.UniqueID.RightAlign()
	suid := tbs.SubjectUniqueID.RightAlign()
	return iuid, suid, err
}

// printName prints the fields of a distinguished name, which include such
// things as its common name and locality.
func printName(names []pkix.AttributeTypeAndValue, buf *bytes.Buffer) []string {
	values := []string{}
	for _, name := range names {
		oid := name.Type
		if len(oid) == 4 && oid[0] == 2 && oid[1] == 5 && oid[2] == 4 {
			switch oid[3] {
			case 3:
				values = append(values, fmt.Sprintf("CN=%s", name.Value))
			case 6:
				values = append(values, fmt.Sprintf("C=%s", name.Value))
			case 8:
				values = append(values, fmt.Sprintf("ST=%s", name.Value))
			case 10:
				values = append(values, fmt.Sprintf("O=%s", name.Value))
			case 11:
				values = append(values, fmt.Sprintf("OU=%s", name.Value))
			default:
				values = append(values, fmt.Sprintf("UnknownOID=%s", name.Type.String()))
			}
		} else if oid.Equal(oidEmailAddress) {
			values = append(values, fmt.Sprintf("emailAddress=%s", name.Value))
		} else {
			values = append(values, fmt.Sprintf("UnknownOID=%s", name.Type.String()))
		}
	}
	if len(values) > 0 {
		buf.WriteString(values[0])
		for i := 1; i < len(values); i++ {
			buf.WriteString("," + values[i])
		}
		buf.WriteString("\n")
	}
	return values
}

// dsaKeyPrinter formats the Y, P, Q, or G components of a DSA public key.
func dsaKeyPrinter(name string, val *big.Int, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%16s%s:", "", name))
	for i, b := range val.Bytes() {
		if (i % 15) == 0 {
			buf.WriteString(fmt.Sprintf("\n%20s", ""))
		}
		buf.WriteString(fmt.Sprintf("%02x", b))
		if i != len(val.Bytes())-1 {
			buf.WriteString(":")
		}
	}
	buf.WriteString("\n")
}

func printVersion(version int, buf *bytes.Buffer) {
	hexVersion := version - 1
	if hexVersion < 0 {
		hexVersion = 0
	}
	buf.WriteString(fmt.Sprintf("%8sVersion: %d (%#x)\n", "", version, hexVersion))
}

func printSubjectInformation(subj *pkix.Name, pkAlgo x509.PublicKeyAlgorithm, pk interface{}, buf *bytes.Buffer) error {
	buf.WriteString(fmt.Sprintf("%8sSubject: ", ""))
	printName(subj.Names, buf)
	buf.WriteString(fmt.Sprintf("%8sSubject Public Key Info:\n%12sPublic Key Algorithm: ", "", ""))
	switch pkAlgo {
	case x509.RSA:
		buf.WriteString(fmt.Sprintf("RSA\n"))
		if rsaKey, ok := pk.(*rsa.PublicKey); ok {
			buf.WriteString(fmt.Sprintf("%16sPublic-Key: (%d bit)\n", "", rsaKey.N.BitLen()))
			// Some implementations (notably OpenSSL) prepend 0x00 to the modulus
			// if its most-significant bit is set. There is no need to do that here
			// because the modulus is always unsigned and the extra byte can be
			// confusing given the bit length.
			buf.WriteString(fmt.Sprintf("%16sModulus:", ""))
			for i, val := range rsaKey.N.Bytes() {
				if (i % 15) == 0 {
					buf.WriteString(fmt.Sprintf("\n%20s", ""))
				}
				buf.WriteString(fmt.Sprintf("%02x", val))
				if i != len(rsaKey.N.Bytes())-1 {
					buf.WriteString(":")
				}
			}
			buf.WriteString(fmt.Sprintf("\n%16sExponent: %d (%#x)\n", "", rsaKey.E, rsaKey.E))
		} else {
			return errors.New("certinfo: Expected rsa.PublicKey for type x509.RSA")
		}
	case x509.DSA:
		buf.WriteString(fmt.Sprintf("DSA\n"))
		if dsaKey, ok := pk.(*dsa.PublicKey); ok {
			dsaKeyPrinter("pub", dsaKey.Y, buf)
			dsaKeyPrinter("P", dsaKey.P, buf)
			dsaKeyPrinter("Q", dsaKey.Q, buf)
			dsaKeyPrinter("G", dsaKey.G, buf)
		} else {
			return errors.New("certinfo: Expected dsa.PublicKey for type x509.DSA")
		}
	case x509.ECDSA:
		buf.WriteString(fmt.Sprintf("ECDSA\n"))
		if ecdsaKey, ok := pk.(*ecdsa.PublicKey); ok {
			buf.WriteString(fmt.Sprintf("%16sPublic-Key: (%d bit)\n", "", ecdsaKey.Params().BitSize))
			dsaKeyPrinter("X", ecdsaKey.X, buf)
			dsaKeyPrinter("Y", ecdsaKey.Y, buf)
			buf.WriteString(fmt.Sprintf("%16sCurve: %s\n", "", ecdsaKey.Params().Name))
		} else {
			return errors.New("certinfo: Expected ecdsa.PublicKey for type x509.DSA")
		}
	default:
		return errors.New("certinfo: Unknown public key type")
	}
	return nil
}

func printSubjKeyId(ext pkix.Extension, buf *bytes.Buffer) error {
	// subjectKeyIdentifier: RFC 5280, 4.2.1.2
	buf.WriteString(fmt.Sprintf("%12sX509v3 Subject Key Identifier:", ""))
	if ext.Critical {
		buf.WriteString(" critical\n")
	} else {
		buf.WriteString("\n")
	}
	var subjectKeyId []byte
	if _, err := asn1.Unmarshal(ext.Value, &subjectKeyId); err != nil {
		return err
	}
	for i := 0; i < len(subjectKeyId); i++ {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("%16s%02X", "", subjectKeyId[0]))
		} else {
			buf.WriteString(fmt.Sprintf(":%02X", subjectKeyId[i]))
		}
	}
	buf.WriteString("\n")
	return nil
}

func printSubjAltNames(ext pkix.Extension, dnsNames []string, emailAddresses []string, ipAddresses []net.IP, buf *bytes.Buffer) error {
	// subjectAltName: RFC 5280, 4.2.1.6
	// TODO: Currently crypto/x509 only extracts DNS, email, and IP addresses.
	// We should add the others to it or implement them here.
	buf.WriteString(fmt.Sprintf("%12sX509v3 Subject Alternative Name:", ""))
	if ext.Critical {
		buf.WriteString(" critical\n")
	} else {
		buf.WriteString("\n")
	}
	if len(dnsNames) > 0 {
		buf.WriteString(fmt.Sprintf("%16sDNS:%s", "", dnsNames[0]))
		for i := 1; i < len(dnsNames); i++ {
			buf.WriteString(fmt.Sprintf(", DNS:%s", dnsNames[i]))
		}
		buf.WriteString("\n")
	}
	if len(emailAddresses) > 0 {
		buf.WriteString(fmt.Sprintf("%16semail:%s", "", emailAddresses[0]))
		for i := 1; i < len(emailAddresses); i++ {
			buf.WriteString(fmt.Sprintf(", email:%s", emailAddresses[i]))
		}
		buf.WriteString("\n")
	}
	if len(ipAddresses) > 0 {
		buf.WriteString(fmt.Sprintf("%16sIP Address:%s", "", ipAddresses[0].String())) // XXX verify string format
		for i := 1; i < len(ipAddresses); i++ {
			buf.WriteString(fmt.Sprintf(", IP Address:%s", ipAddresses[i].String()))
		}
		buf.WriteString("\n")
	}
	return nil
}

func printSignature(sigAlgo x509.SignatureAlgorithm, sig []byte, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%4sSignature Algorithm: %s", "", sigAlgo))
	for i, val := range sig {
		if (i % 18) == 0 {
			buf.WriteString(fmt.Sprintf("\n%9s", ""))
		}
		buf.WriteString(fmt.Sprintf("%02x", val))
		if i != len(sig)-1 {
			buf.WriteString(":")
		}
	}
	buf.WriteString("\n")
}

// CertificateText returns a human-readable string representation
// of the certificate cert. The format is similar (but not identical)
// to the OpenSSL way of printing certificates.
func CertificateText(cert *x509.Certificate) (string, error) {
	var buf bytes.Buffer
	buf.Grow(4096) // 4KiB should be enough

	buf.WriteString(fmt.Sprintf("Certificate:\n"))
	buf.WriteString(fmt.Sprintf("%4sData:\n", ""))
	printVersion(cert.Version, &buf)
	buf.WriteString(fmt.Sprintf("%8sSerial Number: %d (%#x)\n", "", cert.SerialNumber, cert.SerialNumber))
	buf.WriteString(fmt.Sprintf("%4sSignature Algorithm: %s\n", "", cert.SignatureAlgorithm))

	// Issuer information
	buf.WriteString(fmt.Sprintf("%8sIssuer: ", ""))
	printName(cert.Issuer.Names, &buf)

	// Validity information
	buf.WriteString(fmt.Sprintf("%8sValidity\n", ""))
	buf.WriteString(fmt.Sprintf("%12sNot Before: %s\n", "", cert.NotBefore.Format("Jan 2 15:04:05 2006 MST")))
	buf.WriteString(fmt.Sprintf("%12sNot After : %s\n", "", cert.NotAfter.Format("Jan 2 15:04:05 2006 MST")))

	// Subject information
	err := printSubjectInformation(&cert.Subject, cert.PublicKeyAlgorithm, cert.PublicKey, &buf)
	if err != nil {
		return "", err
	}

	// Issuer/Subject Unique ID, typically used in old v2 certificates
	issuerUID, subjectUID, err := certUniqueIDs(cert.RawTBSCertificate)
	if err != nil {
		return "", errors.New(fmt.Sprintf("certinfo: Error parsing TBS unique attributes: %s\n", err.Error()))
	}
	if len(issuerUID) > 0 {
		buf.WriteString(fmt.Sprintf("%8sIssuer Unique ID: %02x", "", issuerUID[0]))
		for i := 1; i < len(issuerUID); i++ {
			buf.WriteString(fmt.Sprintf(":%02x", issuerUID[i]))
		}
		buf.WriteString("\n")
	}
	if len(subjectUID) > 0 {
		buf.WriteString(fmt.Sprintf("%8sSubject Unique ID: %02x", "", subjectUID[0]))
		for i := 1; i < len(subjectUID); i++ {
			buf.WriteString(fmt.Sprintf(":%02x", subjectUID[i]))
		}
		buf.WriteString("\n")
	}

	// Optional extensions for X509v3
	if cert.Version == 3 && len(cert.Extensions) > 0 {
		buf.WriteString(fmt.Sprintf("%8sX509v3 extensions:\n", ""))
		for _, ext := range cert.Extensions {
			if len(ext.Id) == 4 && ext.Id[0] == 2 && ext.Id[1] == 5 && ext.Id[2] == 29 {
				switch ext.Id[3] {
				case 14:
					err = printSubjKeyId(ext, &buf)
				case 15:
					// keyUsage: RFC 5280, 4.2.1.3
					buf.WriteString(fmt.Sprintf("%12sX509v3 Key Usage:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					usages := []string{}
					if cert.KeyUsage&x509.KeyUsageDigitalSignature > 0 {
						usages = append(usages, "Digital Signature")
					}
					if cert.KeyUsage&x509.KeyUsageContentCommitment > 0 {
						usages = append(usages, "Content Commitment")
					}
					if cert.KeyUsage&x509.KeyUsageKeyEncipherment > 0 {
						usages = append(usages, "Key Encipherment")
					}
					if cert.KeyUsage&x509.KeyUsageDataEncipherment > 0 {
						usages = append(usages, "Data Encipherment")
					}
					if cert.KeyUsage&x509.KeyUsageKeyAgreement > 0 {
						usages = append(usages, "Key Agreement")
					}
					if cert.KeyUsage&x509.KeyUsageCertSign > 0 {
						usages = append(usages, "Certificate Sign")
					}
					if cert.KeyUsage&x509.KeyUsageCRLSign > 0 {
						usages = append(usages, "CRL Sign")
					}
					if cert.KeyUsage&x509.KeyUsageEncipherOnly > 0 {
						usages = append(usages, "Encipher Only")
					}
					if cert.KeyUsage&x509.KeyUsageDecipherOnly > 0 {
						usages = append(usages, "Decipher Only")
					}
					if len(usages) > 0 {
						buf.WriteString(fmt.Sprintf("%16s%s", "", usages[0]))
						for i := 1; i < len(usages); i++ {
							buf.WriteString(fmt.Sprintf(", %s", usages[i]))
						}
						buf.WriteString("\n")
					} else {
						buf.WriteString(fmt.Sprintf("%16sNone\n", ""))
					}
				case 17:
					err = printSubjAltNames(ext, cert.DNSNames, cert.EmailAddresses, cert.IPAddresses, &buf)
				case 19:
					// basicConstraints: RFC 5280, 4.2.1.9
					if !cert.BasicConstraintsValid {
						break
					}
					buf.WriteString(fmt.Sprintf("%12sX509v3 Basic Constraints:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					if cert.IsCA {
						buf.WriteString(fmt.Sprintf("%16sCA:TRUE", ""))
					} else {
						buf.WriteString(fmt.Sprintf("%16sCA:FALSE", ""))
					}
					if cert.MaxPathLenZero {
						buf.WriteString(fmt.Sprintf(", pathlen:0\n"))
					} else if cert.MaxPathLen > 0 {
						buf.WriteString(fmt.Sprintf(", pathlen:%d\n", cert.MaxPathLen))
					} else {
						buf.WriteString("\n")
					}
				case 30:
					// nameConstraints: RFC 5280, 4.2.1.10
					// TODO: Currently crypto/x509 only supports "Permitted" and not "Excluded"
					// subtrees. Furthermore it assumes all types are DNS names which is not
					// necessarily true. This missing functionality should be implemented.
					buf.WriteString(fmt.Sprintf("%12sX509v3 Name Constraints:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					if len(cert.PermittedDNSDomains) > 0 {
						buf.WriteString(fmt.Sprintf("%16sPermitted:\n%18s%s", "", "", cert.PermittedDNSDomains[0]))
						for i := 1; i < len(cert.PermittedDNSDomains); i++ {
							buf.WriteString(fmt.Sprintf(", %s", cert.PermittedDNSDomains[i]))
						}
						buf.WriteString("\n")
					}
				case 31:
					// CRLDistributionPoints: RFC 5280, 4.2.1.13
					// TODO: Currently crypto/x509 does not fully implement this section,
					// including types and reason flags.
					buf.WriteString(fmt.Sprintf("%12sX509v3 CRL Distribution Points:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					if len(cert.CRLDistributionPoints) > 0 {
						buf.WriteString(fmt.Sprintf("\n%16sFull Name:\n%18sURI:%s", "", "", cert.CRLDistributionPoints[0]))
						for i := 1; i < len(cert.CRLDistributionPoints); i++ {
							buf.WriteString(fmt.Sprintf(", URI:%s", cert.CRLDistributionPoints[i]))
						}
						buf.WriteString("\n\n")
					}
				case 32:
					// certificatePoliciesExt: RFC 5280, 4.2.1.4
					// TODO: Currently crypto/x509 does not fully impelment this section,
					// including the Certification Practice Statement (CPS)
					buf.WriteString(fmt.Sprintf("%12sX509v3 Certificate Policies:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					for _, val := range cert.PolicyIdentifiers {
						buf.WriteString(fmt.Sprintf("%16sPolicy: %s\n", "", val.String()))
					}
				case 35:
					// authorityKeyIdentifier: RFC 5280, 4.2.1.1
					buf.WriteString(fmt.Sprintf("%12sX509v3 Authority Key Identifier:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					buf.WriteString(fmt.Sprintf("%16skeyid", ""))
					for _, val := range cert.AuthorityKeyId {
						buf.WriteString(fmt.Sprintf(":%02X", val))
					}
					buf.WriteString("\n")
				case 37:
					// extKeyUsage: RFC 5280, 4.2.1.12
					buf.WriteString(fmt.Sprintf("%12sX509v3 Extended Key Usage:", ""))
					if ext.Critical {
						buf.WriteString(" critical\n")
					} else {
						buf.WriteString("\n")
					}
					var list []string
					for _, val := range cert.ExtKeyUsage {
						switch val {
						case x509.ExtKeyUsageAny:
							list = append(list, "Any Usage")
						case x509.ExtKeyUsageServerAuth:
							list = append(list, "TLS Web Server Authentication")
						case x509.ExtKeyUsageClientAuth:
							list = append(list, "TLS Web Client Authentication")
						case x509.ExtKeyUsageCodeSigning:
							list = append(list, "Code Signing")
						case x509.ExtKeyUsageEmailProtection:
							list = append(list, "E-mail Protection")
						case x509.ExtKeyUsageIPSECEndSystem:
							list = append(list, "IPSec End System")
						case x509.ExtKeyUsageIPSECTunnel:
							list = append(list, "IPSec Tunnel")
						case x509.ExtKeyUsageIPSECUser:
							list = append(list, "IPSec User")
						case x509.ExtKeyUsageTimeStamping:
							list = append(list, "Time Stamping")
						case x509.ExtKeyUsageOCSPSigning:
							list = append(list, "OCSP Signing")
						default:
							list = append(list, "UNKNOWN")
						}
					}
					if len(list) > 0 {
						buf.WriteString(fmt.Sprintf("%16s%s", "", list[0]))
						for i := 1; i < len(list); i++ {
							buf.WriteString(fmt.Sprintf(", %s", list[i]))
						}
						buf.WriteString("\n")
					}
				default:
					buf.WriteString(fmt.Sprintf("Unknown extension 2.5.29.%d\n", ext.Id[3]))
				}
				if err != nil {
					return "", err
				}
			} else if ext.Id.Equal(oidExtensionAuthorityInfoAccess) {
				// authorityInfoAccess: RFC 5280, 4.2.2.1
				buf.WriteString(fmt.Sprintf("%12sAuthority Information Access:", ""))
				if ext.Critical {
					buf.WriteString(" critical\n")
				} else {
					buf.WriteString("\n")
				}
				if len(cert.OCSPServer) > 0 {
					buf.WriteString(fmt.Sprintf("%16sOCSP - URI:%s", "", cert.OCSPServer[0]))
					for i := 1; i < len(cert.OCSPServer); i++ {
						buf.WriteString(fmt.Sprintf(",URI:%s", cert.OCSPServer[i]))
					}
					buf.WriteString("\n")
				}
				if len(cert.IssuingCertificateURL) > 0 {
					buf.WriteString(fmt.Sprintf("%16sCA Issuers - URI:%s", "", cert.IssuingCertificateURL[0]))
					for i := 1; i < len(cert.IssuingCertificateURL); i++ {
						buf.WriteString(fmt.Sprintf(",URI:%s", cert.IssuingCertificateURL[i]))
					}
					buf.WriteString("\n")
				}
				buf.WriteString("\n")
			} else if ext.Id.Equal(oidNSComment) {
				// Netscape comment
				var comment string
				rest, err := asn1.Unmarshal(ext.Value, &comment)
				if err != nil || len(rest) > 0 {
					return "", errors.New("certinfo: Error parsing OID " + ext.Id.String())
				}
				if ext.Critical {
					buf.WriteString(fmt.Sprintf("%12sNetscape Comment: critical\n%16s%s\n", "", "", comment))
				} else {
					buf.WriteString(fmt.Sprintf("%12sNetscape Comment:\n%16s%s\n", "", "", comment))
				}
			} else {
				buf.WriteString(fmt.Sprintf("%12sUnknown extension %s\n", "", ext.Id.String()))
			}
		}
		buf.WriteString("\n")
	}

	// Signature
	printSignature(cert.SignatureAlgorithm, cert.Signature, &buf)

	// Optional: Print the full PEM certificate
	/*
		pemBlock := pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}
		buf.Write(pem.EncodeToMemory(&pemBlock))
	*/

	return buf.String(), nil
}

// CertificateRequestText returns a human-readable string representation
// of the certificate request csr. The format is similar (but not identical)
// to the OpenSSL way of printing certificates.
func CertificateRequestText(csr *x509.CertificateRequest) (string, error) {
	var buf bytes.Buffer
	buf.Grow(4096) // 4KiB should be enough

	buf.WriteString(fmt.Sprintf("Certificate Request:\n"))
	buf.WriteString(fmt.Sprintf("%4sData:\n", ""))
	printVersion(csr.Version, &buf)

	// Subject information
	err := printSubjectInformation(&csr.Subject, csr.PublicKeyAlgorithm, csr.PublicKey, &buf)
	if err != nil {
		return "", err
	}

	// Optional extensions for X509v3
	if csr.Version == 3 && len(csr.Extensions) > 0 {
		buf.WriteString(fmt.Sprintf("%8sRequested Extensions:\n", ""))
		var err error
		for _, ext := range csr.Extensions {
			if len(ext.Id) == 4 && ext.Id[0] == 2 && ext.Id[1] == 5 && ext.Id[2] == 29 {
				switch ext.Id[3] {
				case 14:
					err = printSubjKeyId(ext, &buf)
				case 17:
					err = printSubjAltNames(ext, csr.DNSNames, csr.EmailAddresses, csr.IPAddresses, &buf)
				}
			}
			if err != nil {
				return "", err
			}
		}
		buf.WriteString("\n")
	}

	// Signature
	printSignature(csr.SignatureAlgorithm, csr.Signature, &buf)

	return buf.String(), nil
}
