{ buildGoModule, fetchFromGitHub }:
buildGoModule rec {
  pname = "ory";
  version = "0.2.2";

  src = fetchFromGitHub {
    owner = "ory";
    repo = "cli";
    rev = "v${version}";
    sha256 = "sha256-5N69/Gv4eYLbZNN+sEx+RcFyhGCT0hUxDCje1qrbWiY";
  };
  vendorHash = "sha256-J9jyeLIT+1pFnHOUHrzmblVCJikvY05Sw9zMz5qaDOk";

  doCheck = false;
  installPhase = ''
    install -D $GOPATH/bin/cli $out/bin/ory
  '';
}
