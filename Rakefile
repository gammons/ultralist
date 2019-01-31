desc "Builds ultralist for release"

Envs = [
  {goos: "darwin", arch: "386"},
  {goos: "darwin", arch: "amd64"},
  {goos: "darwin", arch: "arm"},
  {goos: "linux", arch: "386"},
  {goos: "linux", arch: "amd64"},
  {goos: "windows", arch: "386"},
  {goos: "windows", arch: "amd64"}
]

Version = "0.9"

task :build do
  `rm -rf dist/#{Version}`
  Envs.each do |env|
    ENV["GOOS"] = env[:goos]
    ENV["GOARCH"] = env[:arch]
    puts "Building #{env[:goos]} #{env[:arch]}"
    `GOOS=#{env[:goos]} GOARCH=#{env[:arch]} go build -v -o dist/#{Version}/ultralist`
    puts "Tarring #{env[:goos]} #{env[:arch]}"
    `tar -czvf dist/#{Version}/ultralist#{env[:goos]}_#{env[:arch]}.tar.gz dist/#{Version}/ultralist`
    puts "Removing dist/#{Version}/ultralist"
    `rm -rf dist/#{Version}/ultralist`
  end
end

desc "Tests all the things"
task :test do
  system "go test ./..."
end

task default: :test
