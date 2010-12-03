
GOFILES="names.go HomePage.go"

([[ -d _obj ]] || mkdir -p _obj) \
    && 8g -o _obj/dusell.8 $GOFILES \
    && gopack grc _obj/dusell.a _obj/dusell.8 \
    && 8g -o _obj/main.8 main.go \
    && 8l -o main _obj/main.8 \
    && ./main $@
