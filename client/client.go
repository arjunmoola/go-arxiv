package client

import (
	"encoding/xml"
	"context"
	"net/http"
	"net/url"
	"log"
	"fmt"
	"io"
)

type SearchOp string

const (
	AND SearchOp = "AND"
	OR SearchOp = "OR"
	ANDNOT SearchOp = "ANDNOT"
)

type SearchSortOrder string

const (
	ASCENDING SearchSortOrder = "ascending"
	DESCENDING SearchSortOrder = "descending"
)

type SearchSortBy string

const (
	RELEVANCE SearchSortBy = "relevance"
	LASTUPDATED SearchSortBy = "lastUpdatedDate"
	SUBMITTED SearchSortBy = "submittedDate"
)

type FieldPrefix int

const (
	Title FieldPrefix = iota
	AuthorPrefix
	Abstract
	Comment
	JournalReference
	SubjectCategory
	ReportNumber
	Id
	AllOfTheAbove
)

func (f FieldPrefix) Code() string {
	var code string
	switch f {
	case Title:
		code = "ti"
	case AuthorPrefix:
		code = "au"
	case Abstract:
		code = "abs"
	case Comment:
		code = "co"
	case JournalReference:
		code = "jr"
	case SubjectCategory:
		code = "cat"
	case ReportNumber:
		code = "rn"
	case Id:
		code = "id"
	case AllOfTheAbove:
		code = "all"
	default:
		code =  "unk"
	}

	return code
}

const arxivUrlStr = "http://export.arxiv.org/api/query"

type urlConstructor struct {
	u *url.URL
	value url.Values
}

func newUrlConstructor(base string) (*urlConstructor, error) {
	u, err := url.Parse(base)

	if err != nil {
		return nil, err
	}

	constructor := &urlConstructor{
		u: u,
		value: make(url.Values),
	}

	return constructor, nil
}

type SearchOperator func(u *urlConstructor)


func WithAnd(prefix FieldPrefix, q string) SearchOperator {
	return func(u *urlConstructor) {
		updateSearchQuery(u, AND, prefix, q)
	}
}

func WithOr(prefix FieldPrefix, q string) SearchOperator {
	return func(u *urlConstructor) {
		updateSearchQuery(u, OR, prefix, q)
	}
}

func WithAndNot(prefix FieldPrefix, q string) SearchOperator {
	return func(u *urlConstructor) {
		updateSearchQuery(u, ANDNOT, prefix, q)
	}
}

func WithMaxResults(n int) SearchOperator {
	return func(u *urlConstructor) {
		u.value.Set("max_results", fmt.Sprintf("%d", n))
	}
}

func WithSortby(value SearchSortBy) SearchOperator {
	return func(u *urlConstructor) {
		u.value.Set("sortBy", string(value))
	}
}

func WithSortOrder(value SearchSortOrder) SearchOperator {
	return func(u *urlConstructor) {
		u.value.Set("sortOrder", string(value))
	}
}

func updateSearchQuery(u *urlConstructor, op SearchOp, prefix FieldPrefix, q string) {
	val := u.value.Get("search_query")
	if val == "" {
		u.value.Add("search_query", prefix.Code()+":"+q)
		return
	}
	val += " " + string(op) + " " + prefix.Code()+":"+q
	u.value.Set("search_query", val)
}

func (u *urlConstructor) searchQuery(prefix FieldPrefix, q string, ops ...SearchOperator) {
	u.value.Set("search_query", prefix.Code()+":"+q)

	for _, op := range ops {
		op(u)
	}
}

func (u *urlConstructor) String() string {
	u.u.RawQuery = u.value.Encode()

	return u.u.String()
}

func (u *urlConstructor) NewRequestWithCtx(ctx context.Context, method string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, u.String(), body)
}

// Feed represents the root element of the arXiv API Atom feed response
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Links   []Link   `xml:"link"`
	Entries []Entry  `xml:"entry"`
	
	// OpenSearch elements (namespace prefix handled via xml.Name)
	TotalResults OpenSearchInt `xml:"http://a9.com/-/spec/opensearch/1.1/ totalResults"`
	StartIndex   OpenSearchInt `xml:"http://a9.com/-/spec/opensearch/1.1/ startIndex"`
	ItemsPerPage OpenSearchInt `xml:"http://a9.com/-/spec/opensearch/1.1/ itemsPerPage"`
}

// OpenSearchInt is a helper type for OpenSearch integer elements
type OpenSearchInt struct {
	Value int `xml:",chardata"`
}

// Link represents a link element in the Atom feed
type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	Title   string   `xml:"title,attr"`
}

// Entry represents a single paper entry in the arXiv feed
type Entry struct {
	XMLName xml.Name `xml:"entry"`
	
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Published string `xml:"published"`
	Title   string   `xml:"title"`
	Summary string   `xml:"summary"`
	Authors []Author `xml:"author"`
	Links   []Link   `xml:"link"`
	
	// arXiv-specific elements (using full namespace URI)
	PrimaryCategory Category `xml:"http://arxiv.org/schemas/atom primary_category"`
	Categories      []Category `xml:"http://arxiv.org/schemas/atom category"`
	Comment         string     `xml:"http://arxiv.org/schemas/atom comment,omitempty"`
	JournalRef      string     `xml:"http://arxiv.org/schemas/atom journal_ref,omitempty"`
	DOI             string     `xml:"http://arxiv.org/schemas/atom doi,omitempty"`
}

// Author represents an author of a paper
type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Affiliation string `xml:"http://arxiv.org/schemas/atom affiliation,omitempty"`
}

// Category represents a subject category
type Category struct {
	Term    string   `xml:"term,attr"`
	Scheme  string   `xml:"scheme,attr"`
}

type Client struct {
	client *http.Client
}

func New() *Client {
	return &Client{
		client: &http.Client{},
	}
}

func (c *Client) Search(ctx context.Context, prefix FieldPrefix, query string, ops ...SearchOperator) (Feed, error) {
	var feed Feed
	
	url, err := newUrlConstructor(arxivUrlStr)

	if err != nil {
		return feed, err
	}

	url.searchQuery(prefix, query, ops...)

	req, err := url.NewRequestWithCtx(ctx, http.MethodGet, nil)

	if err != nil {
		return feed, err
	}

	if err := fetchRequest(c, req, &feed); err != nil {
		return feed, err
	}

	return feed, nil
}

func fetchRequest(c *Client, req *http.Request, v any) error {
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return xml.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
