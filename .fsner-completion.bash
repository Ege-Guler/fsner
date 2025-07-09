#!/bin/bash

__fsner_bash_autocomplete() {
  local cur opts words
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  opts=$(fsner --generate-bash-completion 2>/dev/null)
  COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
}

complete -F __fsner_bash_autocomplete fsner
