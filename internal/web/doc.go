// Package web contains the web-app that sits in front of the
// registry server to provide a nice UI for end users.  This is a
// react app which has to be built for statik inclusion.
package web

//go:generate yarn install
//go:generate yarn run build
//go:generate fileb0x fileb0x.yml
