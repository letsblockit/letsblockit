{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-3QCQqjcCPoNz1GU4SVbH+DERuq6pxwlVu5eW9KaJ/ew=";
  version = "1.0";
}
