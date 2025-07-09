__fsner_bash_autocomplete() {
  local cur opts words
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  words=("${COMP_WORDS[@]:0:$COMP_CWORD}")
  requestComp="${words[*]} --generate-bash-completion"
  opts=$(eval "${requestComp}" 2>/dev/null)
  COMPREPLY=( $(compgen -W "${opts}" -- "$cur") )
}
complete -o bashdefault -o default -o nospace -o nosort -F __fsner_bash_autocomplete fsner
