cd sql/schema &&
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down-to 0 && goose postgres "postgres://postgres:postgres@localhost:5432/gator" up &&
cd ../..

