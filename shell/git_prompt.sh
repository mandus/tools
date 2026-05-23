# Git branch prompt for BusyBox ash (POSIX sh compatible)
# Requires: git >= 2.13 (for --porcelain=v2)
#
# Performance: 1 git call (tags disabled) or 2 git calls (tags enabled).
# Previous version used 5 git subprocess calls per prompt.
#
# Environment variables:
#   GIT_PROMPT_DISABLE_TAGS=1  — skip tag lookup (saves one git fork)
#
# Source this file in your ~/.profile, BEFORE zoxide init:
#   . /path/to/git_prompt.sh
#   eval "$(zoxide init posix --hook prompt)"

__git_branch() {
    _gp_status=$(git status --porcelain=v2 --branch 2>/dev/null) || return
    _gp_branch=
    _gp_tag=
    _gp_sync=
    _gp_ahead=0
    _gp_behind=0
    _gp_has_ab=
    _gp_initial=
    _gp_oldifs=$IFS
    _gp_nl='
'
    IFS=$_gp_nl
    for _gp_line in $_gp_status; do
        case $_gp_line in
            "# branch.head "*)
                _gp_branch=${_gp_line#\# branch.head }
                ;;
            "# branch.oid (initial)")
                _gp_initial=1
                ;;
            "# branch.ab "*)
                _gp_has_ab=1
                _gp_line=${_gp_line#\# branch.ab }
                _gp_ahead=${_gp_line%% *}
                _gp_behind=${_gp_line#* }
                _gp_ahead=${_gp_ahead#+}
                _gp_behind=${_gp_behind#-}
                ;;
        esac
    done
    IFS=$_gp_oldifs

    [ -n "$_gp_branch" ] || return

    # Match old behavior: detached HEAD and initial repo both show "HEAD"
    if [ "$_gp_branch" = "(detached)" ] || [ "$_gp_initial" = "1" ]; then
        _gp_branch=HEAD
    fi

    if [ "$_gp_has_ab" = "1" ]; then
        if [ "$_gp_ahead" -gt 0 ] && [ "$_gp_behind" -gt 0 ]; then
            _gp_sync="<>"
        elif [ "$_gp_ahead" -gt 0 ]; then
            _gp_sync=">"
        elif [ "$_gp_behind" -gt 0 ]; then
            _gp_sync="<"
        else
            _gp_sync="="
        fi
    fi

    if [ "${GIT_PROMPT_DISABLE_TAGS:-}" != "1" ]; then
        _gp_tag=$(git describe --tags --abbrev=0 2>/dev/null)
        if [ -n "$_gp_tag" ]; then
            _gp_tag="($_gp_tag)"
        fi
    fi

    printf '\033[;32m(%s%s%s)\033[0m' "$_gp_branch" "$_gp_tag" "$_gp_sync"
}

# Single prompt function — avoids multiple $() in PS1 which
# can be unreliable in some BusyBox ash builds.
__prompt_cmd() {
    __git_branch
    \command zoxide add -- "$(__zoxide_pwd 2>/dev/null || \command pwd -P)" 2>/dev/null
}

#Add this in .profile after eval zoxide
#Also set _ZO_DOCTOR=0 to avoid zoxide error message
#export PS1='\033]0;\W\007\033[0;33m\w\033[0m$(__prompt_cmd)|\033[;34m\t\033[0m>'

# Interactive git checkout with fzf branch selection
# Usage: gco [partial]   — fuzzy-select a branch to checkout
#        gco branch_name — checkout directly if exact match
gco() {
    _gco_branches() {
        #git branch --format='%(refname:short)' 2>/dev/null
        git branch -r --format='%(refname:short)' 2>/dev/null \
            | grep -v '^origin$' | sed 's|^[^/]*/||'
    }
    if [ $# -eq 0 ]; then
        branch=$(_gco_branches | sort -u \
            | fzf --height=40% --reverse --prompt="checkout> ")
    else
        # If exact match exists, checkout directly
        if git show-ref --verify --quiet "refs/heads/$1" 2>/dev/null; then
            branch="$1"
        else
            # Otherwise fuzzy-filter with the argument as query
            branch=$(_gco_branches | sort -u \
                | fzf --height=40% --reverse --prompt="checkout> " --query="$1")
        fi
    fi
    [ -n "$branch" ] && git checkout "$branch"
}
