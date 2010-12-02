([[ -d _obj ]] || mkdir -p _obj) \
    && 8g -o _obj/web.8 src/app.go src/view.go src/cgi.go \
    && gopack grc _obj/web.a _obj/web.8 \
    && 8g -o _obj/test.8 test.go \
    && 8l -o test _obj/test.8 \
    && ./test

#    && gopack grc _obj/web.a _obj/cgi.8 \
