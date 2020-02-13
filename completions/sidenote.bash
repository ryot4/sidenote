_sidenote_path()
{
    local -r opts="$1"
    local cur dir base
    _get_comp_words_by_ref cur

    sidenote ${opts} path -c > /dev/null 2>&1 || return

    case "${cur}" in
    */)
        dir="${cur}"
        ;;
    *)
        base=$(basename "${cur}")
        if [[ ${base} != ${cur} ]]; then
            dir=$(dirname "${cur}")/
        fi
        ;;
    esac
    COMPREPLY=($(compgen -P "${dir}" -W "$(sidenote ${opts} ls "${dir}")" -- "${base}"))
    if [[ ${#COMPREPLY[@]} -gt 1 ]] || [[ ${COMPREPLY[0]} == */ ]]; then
        compopt -o nospace
    fi
}

_sidenote_path_one()
{
    local i=1 nonopts=0
    while [[ $i -lt ${COMP_CWORD} ]]; do
        if [[ ${COMP_WORDS[i]} != -* ]]; then
            nonopts=$((nonopts + 1))
        fi
        i=$((i + 1))
    done
    [[ ${nonopts} -le 1 ]] && _sidenote_path "$@"
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
            _sidenote_path_one "${opts}"
            ;;
        esac
        ;;
    ls)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l -r -t' -- "${cur}"))
            ;;
        *)
            _sidenote_path_one "${opts}"
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
                _sidenote_path_one "${opts}"
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
            _sidenote_path_one "${opts}"
            ;;
        esac
        ;;
    esac
}
complete -F _sidenote sidenote
