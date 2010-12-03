(  [[ -d _obj ]] || mkdir -p _obj ) \
    && 8g -o _obj/duzyinfo.8 names.go HomePage.go \
    && gopack grc _obj/duzyinfo.a _obj/duzyinfo.8 \
    && 8g -o _obj/main.8 main.go \
    && 8l -o main _obj/main.8 \
    && ./main $@
