// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/server/deps.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	news "github.com/letsblockit/letsblockit/src/news"
	pages "github.com/letsblockit/letsblockit/src/pages"
)

// MockPageRenderer is a mock of PageRenderer interface.
type MockPageRenderer struct {
	ctrl     *gomock.Controller
	recorder *MockPageRendererMockRecorder
}

// MockPageRendererMockRecorder is the mock recorder for MockPageRenderer.
type MockPageRendererMockRecorder struct {
	mock *MockPageRenderer
}

// NewMockPageRenderer creates a new mock instance.
func NewMockPageRenderer(ctrl *gomock.Controller) *MockPageRenderer {
	mock := &MockPageRenderer{ctrl: ctrl}
	mock.recorder = &MockPageRendererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPageRenderer) EXPECT() *MockPageRendererMockRecorder {
	return m.recorder
}

// BuildPageContext mocks base method.
func (m *MockPageRenderer) BuildPageContext(c echo.Context, title string) *pages.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildPageContext", c, title)
	ret0, _ := ret[0].(*pages.Context)
	return ret0
}

// BuildPageContext indicates an expected call of BuildPageContext.
func (mr *MockPageRendererMockRecorder) BuildPageContext(c, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildPageContext", reflect.TypeOf((*MockPageRenderer)(nil).BuildPageContext), c, title)
}

// Redirect mocks base method.
func (m *MockPageRenderer) Redirect(c echo.Context, code int, target string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Redirect", c, code, target)
	ret0, _ := ret[0].(error)
	return ret0
}

// Redirect indicates an expected call of Redirect.
func (mr *MockPageRendererMockRecorder) Redirect(c, code, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redirect", reflect.TypeOf((*MockPageRenderer)(nil).Redirect), c, code, target)
}

// RedirectToPage mocks base method.
func (m *MockPageRenderer) RedirectToPage(c echo.Context, name string, params ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{c, name}
	for _, a := range params {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RedirectToPage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RedirectToPage indicates an expected call of RedirectToPage.
func (mr *MockPageRendererMockRecorder) RedirectToPage(c, name interface{}, params ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{c, name}, params...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedirectToPage", reflect.TypeOf((*MockPageRenderer)(nil).RedirectToPage), varargs...)
}

// RegisterContextBuilder mocks base method.
func (m *MockPageRenderer) RegisterContextBuilder(b pages.ContextBuilder) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterContextBuilder", b)
}

// RegisterContextBuilder indicates an expected call of RegisterContextBuilder.
func (mr *MockPageRendererMockRecorder) RegisterContextBuilder(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterContextBuilder", reflect.TypeOf((*MockPageRenderer)(nil).RegisterContextBuilder), b)
}

// RegisterHelpers mocks base method.
func (m *MockPageRenderer) RegisterHelpers(helpers map[string]interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterHelpers", helpers)
}

// RegisterHelpers indicates an expected call of RegisterHelpers.
func (mr *MockPageRendererMockRecorder) RegisterHelpers(helpers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterHelpers", reflect.TypeOf((*MockPageRenderer)(nil).RegisterHelpers), helpers)
}

// Render mocks base method.
func (m *MockPageRenderer) Render(c echo.Context, name string, data *pages.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Render", c, name, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Render indicates an expected call of Render.
func (mr *MockPageRendererMockRecorder) Render(c, name, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*MockPageRenderer)(nil).Render), c, name, data)
}

// RenderWithSidebar mocks base method.
func (m *MockPageRenderer) RenderWithSidebar(c echo.Context, name, sidebar string, data *pages.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderWithSidebar", c, name, sidebar, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenderWithSidebar indicates an expected call of RenderWithSidebar.
func (mr *MockPageRendererMockRecorder) RenderWithSidebar(c, name, sidebar, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderWithSidebar", reflect.TypeOf((*MockPageRenderer)(nil).RenderWithSidebar), c, name, sidebar, data)
}

// MockReleaseClient is a mock of ReleaseClient interface.
type MockReleaseClient struct {
	ctrl     *gomock.Controller
	recorder *MockReleaseClientMockRecorder
}

// MockReleaseClientMockRecorder is the mock recorder for MockReleaseClient.
type MockReleaseClientMockRecorder struct {
	mock *MockReleaseClient
}

// NewMockReleaseClient creates a new mock instance.
func NewMockReleaseClient(ctrl *gomock.Controller) *MockReleaseClient {
	mock := &MockReleaseClient{ctrl: ctrl}
	mock.recorder = &MockReleaseClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReleaseClient) EXPECT() *MockReleaseClientMockRecorder {
	return m.recorder
}

// GetLatestAt mocks base method.
func (m *MockReleaseClient) GetLatestAt() (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestAt")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestAt indicates an expected call of GetLatestAt.
func (mr *MockReleaseClientMockRecorder) GetLatestAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestAt", reflect.TypeOf((*MockReleaseClient)(nil).GetLatestAt))
}

// GetReleases mocks base method.
func (m *MockReleaseClient) GetReleases() ([]*news.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleases")
	ret0, _ := ret[0].([]*news.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleases indicates an expected call of GetReleases.
func (mr *MockReleaseClientMockRecorder) GetReleases() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleases", reflect.TypeOf((*MockReleaseClient)(nil).GetReleases))
}
