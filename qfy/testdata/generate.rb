#!/usr/bin/env ruby
require 'json'

SRC_RANDOM = Random.new(1984)
MAX_VALUES = Hash.new(1).update "rcat" => 50, "vcat" => 50, "infq" => 5, "kws" => 10
MRG_VALUES = {
  "dev"  => ["oth","oth","oth","oth","oth","oth"],
  "bwsm" => ["ff", "sf", "op", "ng", "kq", "an", "ms", "kk", "mo"],
  "pos"  => [2, 4, 5, 6, 7],
  "mob"  => [0] * 8,
  "loc"  => ["en","en","en","en","en","en","en","en","en","en","en","en","es","es","es","es"],
}

targets = []
attrs   = {}

File.open(File.expand_path("../targets.json", __FILE__)) do |file|
  targets = JSON.load(file)
end

targets.each do |target|
  target['rules'].each do |rule|
    attrs[rule['attr']] ||= []
    attrs[rule['attr']].concat rule['values']
  end
end

attrs.keys.each do |name|
  attrs[name] = attrs[name].sort.uniq
  attrs[name].concat(MRG_VALUES[name] || [])
end

File.open(File.expand_path("../facts.json", __FILE__), "w") do |file|
  10000.times do |_|
    fact = {}
    attrs.each do |name, values|
      vals = if MAX_VALUES[name] == 1
        values.sample
      else
        values.sample(SRC_RANDOM.rand(MAX_VALUES[name])+1)
      end
      fact[name] = vals
    end
    file.puts JSON.dump(fact)
  end
end

p attrs.keys
