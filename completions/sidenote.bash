_sidenote_path()
{
    local opts cur dir base
    _get_comp_words_by_ref cur

    if [[ -n "$1" ]]; then
        opts="-d $1"
    fi

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

_sidenote()
{
    local -r cmds=(cat completion edit exec import init ls path rm serve show)
    local cur prev
    _get_comp_words_by_ref cur prev

    local i=1 word cmd notes
    while [[ $i -lt ${COMP_CWORD} ]] && [[ -z ${cmd} ]]; do
        word="${COMP_WORDS[i]}"
        case "${word}" in
        -d)
            notes="${COMP_WORDS[i+1]}"
            __expand_tilde_by_ref notes
            ;;
        *)
            for c in "${cmds[@]}"; do
                if [[ ${c} == ${word} ]]; then
                    cmd="${word}"
                    break
                fi
            done
            ;;
        esac
        i=$((i + 1))
    done

    COMPREPLY=()
    case "${cmd}" in
    '')
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-d -h -V' -- "${cur}"))
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
        ;;
    cat)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${notes}"
            ;;
        esac
        ;;
    completion)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h' -- "${cur}"))
            ;;
        *)
            COMPREPLY=($(compgen -W 'bash' -- "${cur}"))
            ;;
        esac
        ;;
    edit)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-f -h -p' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -f)
                ;;
            *)
                _sidenote_path "${notes}"
                ;;
            esac
            ;;
        esac
        ;;
    exec)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-cd -h' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -cd)
                # XXX: We should complete only directories here.
                _sidenote_path "${opts}"
                ;;
            *)
                # XXX: compgen -c includes shell functions and builtins.
                COMPREPLY=($(compgen -c -- "${cur}"))
                ;;
            esac
            ;;
        esac
        ;;
    import)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-f -d -h' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -)
                _sidenote_path "${notes}"
                ;;
            import|-*)
                _filedir
                ;;
            *)
                _sidenote_path "${notes}"
                ;;
            esac
            ;;
        esac
        ;;
    init)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -l)
                _filedir -d
                ;;
            esac
            ;;
        esac
        ;;
    ls)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l -r -t' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            ls|-*)
                _sidenote_path "${notes}"
                ;;
            esac
            ;;
        esac
        ;;
    path)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-a -c -h -L' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            path|-*)
                _sidenote_path "${notes}"
                ;;
            esac
            ;;
        esac
        ;;
    rm)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -r' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${notes}"
            ;;
        esac
        ;;
    serve)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h -l -t' -- "${cur}"))
            ;;
        *)
            case "${prev}" in
            -l|-t)
                ;;
            *)
                # XXX: We should complete only directories here.
                _sidenote_path "${notes}"
                ;;
            esac
            ;;
        esac
        ;;
    show)
        case "${cur}" in
        -*)
            COMPREPLY=($(compgen -W '-h' -- "${cur}"))
            ;;
        *)
            _sidenote_path "${notes}"
            ;;
        esac
        ;;
    esac
}

complete -F _sidenote sidenote
