module github.com/vvakame/fosite-datastore-storage/example

go 1.12

require (
	cloud.google.com/go v0.38.0
	github.com/favclip/golidator v2.1.0+incompatible
	github.com/favclip/ucon v0.0.0-20190502090340-f7a2801cdedb
	github.com/google/wire v0.2.1
	github.com/ory/fosite v0.29.6
	github.com/vvakame/fosite-datastore-storage/v2 v2.0.0-20181111163114-0e97ec9aa6dd
	go.mercari.io/datastore v1.4.0
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/tools v0.0.0-20190503185657-3b6f9c0030f7
	golang.org/x/xerrors v0.0.0-20190410155217-1f06c39b4373
	honnef.co/go/tools v0.0.0-20190418001031-e561f6794a2a
)

replace github.com/vvakame/fosite-datastore-storage/v2 => ../v2
