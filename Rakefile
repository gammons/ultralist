desc "Builds todolist for release"

Envs = [
  {goos: "darwin", arch: "386"},
  {goos: "darwin", arch: "amd64"},
  {goos: "darwin", arch: "arm"},
  {goos: "linux", arch: "386"},
  {goos: "linux", arch: "amd64"},
  {goos: "windows", arch: "386"},
  {goos: "windows", arch: "amd64"}
]

Version = "0.8.1"

task :build do
  `rm -rf dist/#{Version}`
  Envs.each do |env|
    ENV["GOOS"] = env[:goos]
    ENV["GOARCH"] = env[:arch]
    puts "Building #{env[:goos]} #{env[:arch]}"
    `GOOS=#{env[:goos]} GOARCH=#{env[:arch]} go build -v -o dist/#{Version}/todolist`
    puts "Tarring #{env[:goos]} #{env[:arch]}"
    `tar -czvf dist/#{Version}/todolist_#{env[:goos]}_#{env[:arch]}.tar.gz dist/#{Version}/todolist`
    puts "Removing dist/#{Version}/todolist"
    `rm -rf dist/#{Version}/todolist`
  end
end

desc "Rebuilds vendor dir, assumes only a vendor/vendor.json file in ./vendor"
task :rebuild_vendor do
  system "govendor fetch +missing,^program"
end

desc "Tests all the things"
task :test do
  system "go test ./..."
end

task default: :test
