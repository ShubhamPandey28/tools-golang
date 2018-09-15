// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// ===== Parser creation info state change tests =====
func TestParser2_1CIMovesToPackageAfterParsingPackageNameTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	pkgName := "testPkg"
	err := parser.parsePair2_1("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new package")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != pkgName {
		t.Errorf("expected package name %s, got %s", pkgName, parser.pkg.PackageName)
	}
	// and the package should _not_ be an "unpackaged" placeholder
	if parser.pkg.IsUnpackaged == true {
		t.Errorf("package incorrectly has IsUnpackaged flag set")
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the package should be in the SPDX Document's slice of packages
	flagFound := false
	for _, p := range parser.doc.Packages {
		if p == parser.pkg {
			flagFound = true
		}
	}
	if flagFound == false {
		t.Errorf("package isn't in the SPDX Document's slice of packages")
	}
}

func TestParser2_1CIMovesToFileAfterParsingFileNameTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("FileName", "testFile")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}
	// and current package should be an "unpackaged" placeholder
	if parser.pkg == nil {
		t.Fatalf("parser didn't create placeholder package")
	}
	if !parser.pkg.IsUnpackaged {
		t.Errorf("placeholder package is not set as unpackaged")
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the package should be in the SPDX Document's slice of packages
	flagFound := false
	for _, p := range parser.doc.Packages {
		if p == parser.pkg {
			flagFound = true
		}
	}
	if flagFound == false {
		t.Errorf("package isn't in the SPDX Document's slice of packages")
	}
}

func TestParser2_1CIMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense2_1)
	}
}

func TestParser2_1CIStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}
}

// ===== Creation Info section tests =====
func TestParser2_1HasCreationInfoAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePairFromCreationInfo2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("got error when calling parsePairFromCreationInfo2_1: %v", err)
	}
	if parser.doc.CreationInfo == nil {
		t.Errorf("doc.CreationInfo is still nil after parsing first pair")
	}
}

func TestParser2_1CanParseCreationInfoTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// SPDX Version
	err := parser.parsePairFromCreationInfo2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.SPDXVersion != "SPDX-2.1" {
		t.Errorf("got %v for SPDXVersion", parser.doc.CreationInfo.SPDXVersion)
	}

	// Data License
	err = parser.parsePairFromCreationInfo2_1("DataLicense", "CC0-1.0")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DataLicense != "CC0-1.0" {
		t.Errorf("got %v for DataLicense", parser.doc.CreationInfo.DataLicense)
	}

	// SPDX Identifier
	err = parser.parsePairFromCreationInfo2_1("SPDXID", "SPDXRef-DOCUMENT")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.SPDXIdentifier != "SPDXRef-DOCUMENT" {
		t.Errorf("got %v for SPDXIdentifier", parser.doc.CreationInfo.SPDXIdentifier)
	}

	// Document Name
	err = parser.parsePairFromCreationInfo2_1("DocumentName", "xyz-2.1.5")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentName != "xyz-2.1.5" {
		t.Errorf("got %v for DocumentName", parser.doc.CreationInfo.DocumentName)
	}

	// Document Namespace
	err = parser.parsePairFromCreationInfo2_1("DocumentNamespace", "http://example.com/xyz-2.1.5.spdx")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentNamespace != "http://example.com/xyz-2.1.5.spdx" {
		t.Errorf("got %v for DocumentNamespace", parser.doc.CreationInfo.DocumentNamespace)
	}

	// External Document Reference
	refs := []string{
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1: d6a770ba38583ed4bb4525bd96e50461655d2759",
		"DocumentRef-xyz-2.1.2 http://example.com/xyz-2.1.2 SHA1: d6a770ba38583ed4bb4525bd96e50461655d2760",
	}
	err = parser.parsePairFromCreationInfo2_1("ExternalDocumentRef", refs[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("ExternalDocumentRef", refs[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.ExternalDocumentReferences) != 2 ||
		parser.doc.CreationInfo.ExternalDocumentReferences[0] != refs[0] ||
		parser.doc.CreationInfo.ExternalDocumentReferences[1] != refs[1] {
		t.Errorf("got %v for ExternalDocumentReferences", parser.doc.CreationInfo.ExternalDocumentReferences)
	}

	// License List Version
	err = parser.parsePairFromCreationInfo2_1("LicenseListVersion", "2.2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.LicenseListVersion != "2.2" {
		t.Errorf("got %v for LicenseListVersion", parser.doc.CreationInfo.LicenseListVersion)
	}

	// Creators: Persons
	refPersons := []string{
		"Person: Person A",
		"Person: Person B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorPersons) != 2 ||
		parser.doc.CreationInfo.CreatorPersons[0] != "Person A" ||
		parser.doc.CreationInfo.CreatorPersons[1] != "Person B" {
		t.Errorf("got %v for CreatorPersons", parser.doc.CreationInfo.CreatorPersons)
	}

	// Creators: Organizations
	refOrgs := []string{
		"Organization: Organization A",
		"Organization: Organization B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorOrganizations) != 2 ||
		parser.doc.CreationInfo.CreatorOrganizations[0] != "Organization A" ||
		parser.doc.CreationInfo.CreatorOrganizations[1] != "Organization B" {
		t.Errorf("got %v for CreatorOrganizations", parser.doc.CreationInfo.CreatorOrganizations)
	}

	// Creators: Tools
	refTools := []string{
		"Tool: Tool A",
		"Tool: Tool B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorTools) != 2 ||
		parser.doc.CreationInfo.CreatorTools[0] != "Tool A" ||
		parser.doc.CreationInfo.CreatorTools[1] != "Tool B" {
		t.Errorf("got %v for CreatorTools", parser.doc.CreationInfo.CreatorTools)
	}

	// Created date
	err = parser.parsePairFromCreationInfo2_1("Created", "2018-09-10T11:46:00Z")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.Created != "2018-09-10T11:46:00Z" {
		t.Errorf("got %v for Created", parser.doc.CreationInfo.Created)
	}

	// Creator Comment
	err = parser.parsePairFromCreationInfo2_1("CreatorComment", "Blah whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.CreatorComment != "Blah whatever" {
		t.Errorf("got %v for CreatorComment", parser.doc.CreationInfo.CreatorComment)
	}

	// Document Comment
	err = parser.parsePairFromCreationInfo2_1("DocumentComment", "Blah whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentComment != "Blah whatever" {
		t.Errorf("got %v for DocumentComment", parser.doc.CreationInfo.DocumentComment)
	}

}

func TestParser2_1InvalidCreatorTagsFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("Creator", "blah: somebody")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}

	err = parser.parsePairFromCreationInfo2_1("Creator", "Tool with no colons")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}
}

func TestParser2_1CreatorTagWithMultipleColonsPasses(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("Creator", "Tool: tool1:2:3")
	if err != nil {
		t.Errorf("unexpected error from parsing valid Creator format")
	}
}

func TestParser2_1CIUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParser2_1CICreatesRelationship(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-whatever")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.rln == nil {
		t.Fatalf("parser didn't create and point to Relationship struct")
	}
	if parser.rln != parser.doc.Relationships[0] {
		t.Errorf("pointer to new Relationship doesn't match idx 0 for doc.Relationships[]")
	}
}