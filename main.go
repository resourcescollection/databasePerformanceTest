package main

import "dbtest/databasePerformance"

// postgresql://postgres.qwfejcauesgywwsobofx:GearTest@@123123@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres

// use maybe
//"postgresql://postgres:JPB87TALnMwYa8Mo@irritably-faultless-stonechat.data-1.use1.tembo.io:5432/postgres", "ca.crt"
// europe: perf@skymail.de;
// postgresql://postgres:jQ6YOWxBjaltS2Gy@factually-mature-sponge.data-1.euc1.tembo.io:5432/postgres

// postgresql://872701ec-9694-4792-afe1-7a7b2d769e51-user:pw-c50e83a4-5e5a-431a-a05c-13699c338916@postgres-free-tier-v2020.gigalixir.com:5432/872701ec-9694-4792-afe1-7a7b2d769e51
// warning not suit for product. it's hard convert data from free to product;

// postgresql://mydata_9wqm_user:KVD6zySp9Ux4ou2prT2BNlRo1FNunNse@dpg-csatg1ggph6c73a63s3g-a.singapore-postgres.render.com/mydata_9wqm
// delete in one month

// perf@skymail.de;   // for performance test only;
// postgres://koyeb-adm:SxIFv90GWhro@ep-lively-mode-a1ug026j.ap-southeast-1.pg.koyeb.app/koyebdb

// postgresql://perfs:J78nwp2dVN+tOS73__@sgp.domcloud.co/perfs_db
// PGPASSWORD='J78nwp2dVN+tOS73__' psql -U perfs -h sgp.domcloud.co perfs_db

// postgres://postgres:Sz8Pz7vTCyncnEKE@oqjqekbuuxvdflrndnlh.db.ap-southeast-1.nhost.run:5432/oqjqekbuuxvdflrndnlh

// vercel: github: performance@pissmail.com;
// psql "postgresql://default:p2wODeKGr7Fj@ep-summer-wind-a1a04sob.ap-southeast-1.aws.neon.tech:5432/verceldb"

// filess.io: perf@skymail.de;
// postgresql://test1_pianolamp:4d8213e640b288ea751ef6728f1fcbaba9e8452d@5se8o.h.filess.io:5433/test1_pianolamp

// alwaysdata.com: perf@skymail.de;
// postgresql://perf_u1:GearTest@@123123@postgresql-perf.alwaysdata.net/perf_d1

// "sslmode=disable"
func main() {
	databasePerformance.DP_TestPostgresql("postgresql://koyeb-adm:SxIFv90GWhro@ep-lively-mode-a1ug026j.ap-southeast-1.pg.koyeb.app/koyebdb",
		"")
}
