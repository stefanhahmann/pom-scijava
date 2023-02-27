#!/bin/sh

# Fail on error.
set -e

# Generate the mega melt structure.
tests/run.sh -s

# Pull appropriate branches, build artifacts, and install them locally.
adjust() {
  dir=$1
  remote=$2
  branch=$3
  cd "target/mega-melt/melting-pot/$dir"

  # fetch the needed branch
  git remote add upstream "$remote"
  git fetch upstream

  # discard pom.xml hacks
  git checkout pom.xml

  # switch to the needed branch
  git checkout "$branch"

  # reapply pom.xml hacks
  mv pom.xml pom.xml.original &&
  sed -E -e 's_(https?://maven.imagej.net|http://maven.scijava.org)_https://maven.scijava.org_g' \
    -e 's_http://maven.apache.org/xsd_https://maven.apache.org/xsd_g' pom.xml.original > pom.xml ||
    die "Failed to adjust pom.xml"
  perl -0777 -i -pe 's/(<parent>\s*<groupId>org.scijava<\/groupId>\s*<artifactId>pom-scijava<\/artifactId>\s*<version>)[^\n]*/${1}999-mega-melt<\/version>/igs' pom.xml

  # build and install the component
  mvn -Denforcer.skip -DskipTests -Dmaven.test.skip -Dinvoker.skip clean install

  cd - >/dev/null
}
adjust net.imglib2/imglib2 https://github.com/imglib/imglib2 master
adjust net.imglib2/imglib2-cache https://github.com/imglib/imglib2-cache bump-to-imglib2-6.1.0
adjust sc.fiji/bigdataviewer-core https://github.com/bigdataviewer/bigdataviewer-core imglib2-6.1.0
adjust net.imglib2/imglib2-algorithm https://github.com/imglib/imglib2-algorithm bump-to-imglib2-6.1.0
adjust org.janelia.saalfeldlab/n5-imglib2 https://github.com/saalfeldlab/n5-imglib2 bump-to-imglib2-6.1.0
adjust net.preibisch/multiview-reconstruction https://github.com/PreibischLab/multiview-reconstruction bump-to-imglib2-6.1.0
#adjust com.bitplane/imglib2-imaris-bridge https://github.com/imaris/imglib2-imaris-bridge bump-to-imglib2-6.1.0
adjust org.embl.mobie/mobie-io https://github.com/tpietzsch/mobie-io bump-to-imglib2-6.1.0
adjust saalfeldlab/hot-knife https://github.com/saalfeldlab/hot-knife bump-to-imglib2-6.1.0

# Run the mega melt!
cd target/mega-melt/melting-pot
./melt.sh 2>&1 | tee melt.log
