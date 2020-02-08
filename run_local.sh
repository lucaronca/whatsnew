case "$1" in

"generator")  LOCAL=true go run generator/main.go generator/local.go
    ;;
"uploader")  LOCAL=true s3RssFileName=rss.xml go run uploader/main.go uploader/local.go uploader/fileuploader.go uploader/addfiletos3.go
    ;;
*) echo "Signal number $1 is not processed"
   ;;
esac
