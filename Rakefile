desc "Builds ultralist for release"

Envs = [
  { goos: "darwin", arch: "386" },
  { goos: "darwin", arch: "amd64" },
  { goos: "linux", arch: "arm" },
  { goos: "linux", arch: "arm64" },
  { goos: "linux", arch: "386" },
  { goos: "linux", arch: "amd64" },
  { goos: "windows", arch: "386" },
  { goos: "windows", arch: "amd64" }
].freeze

Version = "1.7.0".freeze

task :build do
  `rm -rf dist/#{Version}`
  Envs.each do |env|
    ENV["GOOS"] = env[:goos]
    ENV["GOARCH"] = env[:arch]
    puts "Building #{env[:goos]} #{env[:arch]}"
    `GOOS=#{env[:goos]} GOARCH=#{env[:arch]} go build -v -o dist/#{Version}/ultralist`
    if env[:goos] == "windows"
      puts "Creating windows executable"
      `mv dist/#{Version}/ultralist dist/#{Version}/ultralist.exe`
      `cd dist/#{Version} && zip ultralist_win.zip ultralist.exe`
      puts "Removing windows executable"
      `rm -rf dist/#{Version}/ultralist.exe`
    else
      puts "Tarring #{env[:goos]} #{env[:arch]}"
      `cd dist/#{Version} && tar -czvf ultralist#{env[:goos]}_#{env[:arch]}.tar.gz ultralist`
      puts "Removing dist/#{Version}/ultralist"
      `rm -rf dist/#{Version}/ultralist`
    end
  end
end

desc "Tests all the things"
task :test do
  system "go test ./..."
end

task default: :test
