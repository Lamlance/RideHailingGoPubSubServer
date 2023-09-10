while getopts "o:" flag
do
  case "${flag}" in
    o) build_service=$OPTARG
  esac
done

if [ "$build_service" == "geo" ];then
  go build ./Services/GeoLocationService
elif [ "$build_service" == "sse" ];then
  go build ./Services/SSEService
elif [ "$build_service" == "trip" ];then
  go build -o ../bin/trip.exe
else
  go build ./Services/GeoLocationService
  go build ./Services/SSEService
  go build ./Services/TripService
fi