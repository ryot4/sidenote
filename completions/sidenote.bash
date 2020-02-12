_sidenote_path()
{
    local -r opts="$1"
    local cur prev path
    _get_comp_words_by_ref cur prev

    sidenote ${opts} path -c > /dev/null 2>&1 || return

    compopt -o nospace
    case "${cur}" in
    */)
        path="${cur}"
        ;;
    esac
    COMPREPLY=($(compgen -W "${path}$(sidenote ${opts} ls ${path})" -- "${cur}"))
}

_sidenote()
{
    local -r cmds=(init path ls cat edit rm)
    local cur prev
    _get_comp_words_by_ref cur prev

    local i=1 word cmd opts notes
    while [[ $i -lt ${COMP_CWORD} ]]; do
        word="${COMP_WORDS[i]}"
        case "${word}" in
        -d)
            notes="${COMP_WORDS[i+1]}"
            __expand_tilde_by_ref notes
            opts="${opts} -d ${notes}"
            ;;
        *)
            for c in "${cmds[@]}"; do
                if [[ ${c} == ${word} ]]; then
                    cmd="${word}"
                    break
                fi
            done
            [[ -n ${cmd} ]] && break
            ;;
        esac
        i=$((i + 1))
    done

    COMPREPLY=()

    if [[ -z ${cmd} ]]; then
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-d -h -version' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -d)
                _filedir -d
                ;;
            *)
                COMPREPLY=($(compgen -W "${cmds[*]}" -- "${cur}"))
                ;;
            esac
            ;;
        esac
        return
    fi

    case "${cmd}" in
    init)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l' -- "${cur}"))
            ;;
        *)
            if [[ ${prev} == '-l' ]]; then
                _filedir -d
            fi
            ;;
        esac
        ;;
    path)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-a -c -h' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${opts}"
            ;;
        esac
        ;;
    ls)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l -r -t' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${opts}"
            ;;
        esac
        ;;
    cat)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${opts}"
            ;;
        esac
        ;;
    edit)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-f -h -x' -- "${cur}"))
            ;;
        *)
            if [[ ${prev} != '-f' ]] && [[ ${prev} != '-x' ]]; then
                _sidenote_path "${opts}"
            fi
            ;;
        esac
        ;;
    rm)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -r' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${opts}"
            ;;
        esac
        ;;
    esac
}
complete -F _sidenote sidenote
