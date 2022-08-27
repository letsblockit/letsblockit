{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-jBckynJuXsps8jv5KefI8+nGXwHdQv8edJffxYO7Wtw=";
  version = "1.0";
}
