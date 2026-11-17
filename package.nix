{
  lib,
  buildGoModule,
  version,
  commit ? "unknown",
  date ? "unknown",
}:

let
  buildinfo = "github.com/mika-rafie/alamino/internal/buildinfo";
in
buildGoModule {
  pname = "alamino";
  inherit version;

  src = lib.cleanSource ./.;
  vendorHash = "sha256-maAVBptEtdrGanJHwAPAmuGBorzIMUgK6T+NmIz1kS0=";
  subPackages = [ "cmd/alamino" ];

  env.CGO_ENABLED = 0;
  ldflags = [
    "-X ${buildinfo}.Version=v${version}"
    "-X ${buildinfo}.Commit=${commit}"
    "-X ${buildinfo}.Date=${date}"
  ];

  # The suite needs a real .git directory, /bin/bash, and wall-clock timing
  # the build sandbox does not provide; ci.yml runs it on three OSes instead.
  doCheck = false;

  meta = {
    description = "Local gate that reviews, tests, and fixes every push before it ships";
    homepage = "https://github.com/mika-rafie/alamino";
    license = lib.licenses.mit;
    mainProgram = "alamino";
  };
}
