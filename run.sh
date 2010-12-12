([[ -d _obj ]] || mkdir -p _obj) \
    && 8g -o _obj/web.8 src/app.go src/view.go src/cgi.go src/sm.go src/db.go src/db_mysql.go src/app_cfg.go \
    && gopack grc _obj/web.a _obj/web.8 \
    && 8g -o _obj/test.8 test.go \
    && 8g -o _obj/test_db.8 test_db.go \
    && 8g -o _obj/test_appcfg.8 test_appcfg.go \
    && 8l -o test _obj/test.8 \
    && 8l -o test_db _obj/test_db.8 \
    && 8l -o test_appcfg _obj/test_appcfg.8 \
    && ./test_appcfg $@

#    && gopack grc _obj/web.a _obj/cgi.8 \
