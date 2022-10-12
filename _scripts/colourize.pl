#!/usr/bin/env perl
use strict;
use warnings;
use diagnostics;
use Term::ANSIColor qw(colored);

while (my $line = <>) {
  print STDOUT process_line($line)."\n";
}

exit(0);

sub process_line {
  my ($line) = @_;
  chomp($line);
  if ($line =~ m!^\[(\d+\-\d+\.\d+)\]\s+([A-Z]+)\s+(.+?)\s*$!) {
    my ($datestamp, $level, $message) = ($1, $2, $3);
    my $colour = "white";
    if ($level eq "ERROR") {
      $colour = "bold white on_red";
    } elsif ($level eq "INFO") {
      $colour = "green";
    } elsif ($level eq "DEBUG") {
      $colour = "yellow";
    }
    my $out = "[".colored($datestamp, "blue")."]";
    $out .= " ".colored($level, $colour);
    if ($level eq "DEBUG") {
      $out .= "\t";
      if ($message =~ m!^(.+?)\:(\d+)\s+\[(.+?)\]\s+(.+?)\s*$!) {
        my ($file, $ln, $tag, $info) = ($1, $2, $3, $4);
        $out .= colored($file, "bright_blue");
        $out .= ":".colored($ln, "blue");
        $out .= " [".colored($tag, "bright_blue")."]";
        $out .= " ".colored($info, "bold cyan");
      } else {
        $out .= $message;
      }
    } elsif ($level eq "ERROR") {
      $out .= "\t".colored($message, $colour);
    } elsif ($level eq "INFO") {
      $out .= "\t".colored($message, $colour);
    } else {
      $out .= "\t".$message;
    }
    return $out;
  }
  return $line;
}
