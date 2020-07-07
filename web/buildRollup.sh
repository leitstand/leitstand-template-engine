#npm install rollup --global!/bin/sh
#npm install rollup --global
#rm -rf ../src/main/resources/META-INF/resources
rollup --config
cp src/*.html ../statik
cp src/*.png ../statik
cp src/*.ico ../statik
cp src/*.css ../statik
cp -R src/images ../statik
cp -R src/openapi ../statik