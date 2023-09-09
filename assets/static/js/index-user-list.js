var loadUserList = (function ($) {
    /**
     * verify user value
     * @param username
     */
    function verifyUser(username) {
        var valid = true;
        if (username.trim() === '' || !/^\w+$/.test(username)) {
            valid = false;
        }
        return {
            valid: valid,
            trim: username
        };
    }

    /**
     * verify token value
     * @param token
     */
    function verifyToken(token) {
        var valid = true;
        if (token.trim() === '' || !/^[\w!@#$%^&*()]+$/.test(token)) {
            valid = false;
        }
        return {
            valid: valid,
            trim: token.trim()
        };
    }

    /**
     * verify comment is valid
     * @param comment
     *
     * @return {{valid:boolean, trim:string}}
     */
    function verifyComment(comment) {
        var valid = true;
        if (comment.trim() !== '' && /[\n\t\r]/.test(comment)) {
            valid = false;
        }
        return {
            valid: valid,
            trim: comment.replace(/[\n\t\r]/g, '')
        };
    }

    /**
     * verify ports is valid
     * @param ports
     *
     * @return {{valid:boolean, trim:string}}
     */
    function verifyPorts(ports) {
        var valid = true;
        if (ports.trim() !== '') {
            try {
                ports.split(",").forEach(function (port) {
                    if (/^\s*\d{1,5}\s*$/.test(port)) {
                        if (parseInt(port) < 1 || parseInt(port) > 65535) {
                            valid = false;
                        }
                    } else if (/^\s*\d{1,5}\s*-\s*\d{1,5}\s*$/.test(port)) {
                        var portRange = port.split('-');
                        if (parseInt(portRange[0]) < 1 || parseInt(portRange[0]) > 65535) {
                            valid = false;
                        } else if (parseInt(portRange[1]) < 1 || parseInt(portRange[1]) > 65535) {
                            valid = false;
                        } else if (parseInt(portRange[0]) > parseInt(portRange[1])) {
                            valid = false;
                        }
                    } else {
                        valid = false;
                    }
                    if (valid === false) {
                        throw 'break';
                    }
                });
            } catch (e) {
            }
        }
        return {
            valid: valid,
            trim: ports.replace(/\s/g, '')
        };
    }

    /**
     * verify domains is valid
     * @param domains
     *
     * @return {{valid:boolean, trim:string}}
     */
    function verifyDomains(domains) {
        var valid = true;
        if (domains.trim() !== '') {
            try {
                domains.split(',').forEach(function (domain) {
                    if (!/^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62}){1,3}$/.test(domain.trim())) {
                        valid = false;
                        throw 'break';
                    }
                });
            } catch (e) {
            }
        }
        return {
            valid: valid,
            trim: domains.replace(/\s/g, '')
        };
    }

    /**
     * verify subdomains is valid
     * @param subdomains
     *
     * @return {{valid:boolean, trim:string}}
     */
    function verifySubdomains(subdomains) {
        var valid = true;
        if (subdomains.trim() !== '') {
            try {
                subdomains.split(',').forEach(function (subdomain) {
                    if (!/^[a-zA-z0-9][a-zA-Z0-9-]{0,19}$/.test(subdomain.trim())) {
                        valid = false;
                        throw 'break';
                    }
                });
            } catch (e) {
            }
        }
        return {
            valid: valid,
            trim: subdomains.replace(/\s/g, '')
        };
    }

    /**
     * load i18n language
     * @param lang {{}} language json
     * @param title page title
     */
    function loadUserList(lang, title) {
        $("#title").text(title);
        var html = layui.laytpl($('#userListTemplate').html()).render();
        $('#content').html(html);

        var apiType = {
            Remove: 1,
            Enable: 2,
            Disable: 3
        }

        //set verify rules
        var verifyRules = {
            user: function (value, item) {
                var result = verifyUser(value);
                if (!result.valid) {
                    return lang['UserFormatError'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            },
            token: function (value, item) {
                var result = verifyToken(value);
                if (!result.valid) {
                    return lang['TokenFormatError'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            },
            comment: function (value, item) {
                var result = verifyComment(value);
                if (!result.valid) {
                    return lang['CommentInvalid'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            },
            ports: function (value, item) {
                var result = verifyPorts(value);
                if (!result.valid) {
                    return lang['PortsInvalid'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            },
            domains: function (value, item) {
                var result = verifyDomains(value);
                if (!result.valid) {
                    return lang['DomainsInvalid'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            },
            subdomains: function (value, item) {
                var result = verifySubdomains(value);
                if (!result.valid) {
                    return lang['SubdomainsInvalid'];
                }
                if (item != null) {
                    if (typeof item === "function") {
                        item && item(result.trim);
                    } else {
                        $(item).val(result.trim);
                    }
                }
            }
        };

        layui.form.set({
            verIncludeRequired: true,
            verify: verifyRules
        });

        var $section = $('#content > section');
        layui.table.render({
            elem: '#tokenTable',
            height: $section.height() - $('#searchForm').height() + 8,
            text: {none: lang['EmptyData']},
            url: '/tokens',
            method: 'get',
            where: {},
            dataType: 'json',
            editTrigger: 'dblclick',
            page: navigator.language.indexOf("zh") !== -1,
            toolbar: '#userListToolbarTemplate',
            defaultToolbar: false,
            cols: [[
                {type: 'checkbox'},
                {field: 'user', title: lang['User'], width: 150, sort: true},
                {field: 'token', title: lang['Token'], width: 200, sort: true, edit: true},
                {field: 'comment', title: lang['Notes'], sort: true, edit: 'textarea'},
                {field: 'ports', title: lang['AllowedPorts'], sort: true, edit: 'textarea'},
                {field: 'domains', title: lang['AllowedDomains'], sort: true, edit: 'textarea'},
                {field: 'subdomains', title: lang['AllowedSubdomains'], sort: true, edit: 'textarea'},
                {
                    field: 'status',
                    title: lang['Status'],
                    width: 100,
                    templet: '<span>{{d.status? "' + lang['Enable'] + '":"' + lang['Disable'] + '"}}</span>',
                    sort: true
                },
                {title: lang['Operation'], width: 150, toolbar: '#userListOperationTemplate'}
            ]]
        });

        /**
         * update layui table data
         * @param obj table update obj
         * @param field update field
         * @param trim new value
         */
        function updateTableField(obj, field, trim) {
            var newData = {};
            newData[field] = trim;
            obj.update(newData);
        }

        layui.table.on('edit(tokenTable)', function (obj) {
            var field = obj.field;
            var value = obj.value;
            var oldValue = obj.oldValue;
            var before = $.extend(true, {}, obj.data);
            var after = $.extend(true, {}, obj.data);
            var verifyMsg = false;
            if (field === 'token') {
                verifyMsg = verifyRules.token(value, function (trim) {
                    updateTableField(obj, field, trim)
                });
                if (verifyMsg) {
                    layui.layer.msg(verifyMsg);
                    return obj.reedit();
                }

                before.token = oldValue;
                after.token = value;
            } else if (field === 'comment') {
                verifyMsg = verifyRules.comment(value, function (trim) {
                    updateTableField(obj, field, trim)
                });
                if (verifyMsg) {
                    layui.layer.msg(verifyMsg);
                    return obj.reedit();
                }

                before.comment = oldValue;
                after.comment = value;
            } else if (field === 'ports') {
                verifyMsg = verifyRules.ports(value, function (trim) {
                    updateTableField(obj, field, trim)
                });
                if (verifyMsg) {
                    layui.layer.msg(verifyMsg);
                    return obj.reedit();
                }

                before.ports = oldValue;
                after.ports = value;
            } else if (field === 'domains') {
                verifyMsg = verifyRules.domains(value, function (trim) {
                    updateTableField(obj, field, trim)
                });
                if (verifyMsg) {
                    layui.layer.msg(verifyMsg);
                    return obj.reedit();
                }

                before.domains = oldValue;
                after.domains = value;
            } else if (field === 'subdomains') {
                verifyMsg = verifyRules.subdomains(value, function (trim) {
                    updateTableField(obj, field, trim)
                });
                if (verifyMsg) {
                    layui.layer.msg(verifyMsg);
                    return obj.reedit();
                }

                before.subdomains = oldValue;
                after.subdomains = value;
            }

            update(before, after);
        });

        layui.table.on('toolbar(tokenTable)', function (obj) {
            var id = obj.config.id;
            var checkStatus = layui.table.checkStatus(id);
            switch (obj.event) {
                case 'add':
                    addPopup();
                    break
                case 'remove':
                    batchRemovePopup(checkStatus.data);
                    break
                case 'disable':
                    batchDisablePopup(checkStatus.data);
                    break
                case 'enable':
                    batchEnablePopup(checkStatus.data);
                    break
            }
        });
        layui.table.on('tool(tokenTable)', function (obj) {
            var data = obj.data;
            switch (obj.event) {
                case 'remove':
                    removePopup(data);
                    break;
                case 'disable':
                    disablePopup(data);
                    break;
                case 'enable':
                    enablePopup(data);
                    break
            }
        });

        /**
         * add user popup
         */
        function addPopup() {
            layui.layer.open({
                type: 1,
                title: lang['NewUser'],
                area: ['500px'],
                content: layui.laytpl(document.getElementById('addUserTemplate').innerHTML).render(),
                btn: [lang['Confirm'], lang['Cancel']],
                btn1: function (index) {
                    if (layui.form.validate('#addUserForm')) {
                        add(layui.form.val('addUserForm'), index);
                    }
                },
                btn2: function (index) {
                    layui.layer.close(index);
                }
            });
        }

        /**
         * add user action
         * @param data {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} user data
         * @param index popup index
         */
        function add(data, index) {
            var loading = layui.layer.load();
            $.ajax({
                url: '/add',
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify(data),
                success: function (result) {
                    if (result.success) {
                        reloadTable();
                        layui.layer.close(index);
                        layui.layer.msg(lang['OperateSuccess'], function (index) {
                            layui.layer.close(index);
                        });
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * update user action
         * @param before {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} data before update
         * @param after {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} data after update
         */
        function update(before, after) {
            var loading = layui.layer.load();
            $.ajax({
                url: '/update',
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify({
                    before: before,
                    after: after,
                }),
                success: function (result) {
                    if (result.success) {
                        layui.layer.msg(lang['OperateSuccess']);
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * batch remove user popup
         * @param data {[{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}]} user data list
         */
        function batchRemovePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmRemoveUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Remove, data, index);
            });
        }

        /**
         * batch disable user popup
         * @param data {[{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}]} user data list
         */
        function batchDisablePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmDisableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Disable, data, index);
            });
        }

        /**
         * batch enable user popup
         * @param data {[{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}]} user data list
         */
        function batchEnablePopup(data) {
            if (data.length === 0) {
                layui.layer.msg(lang['ShouldCheckUser']);
                return;
            }
            layui.layer.confirm(lang['ConfirmEnableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Enable, data, index);
            });
        }

        /**
         * remove one user popup
         * @param data {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} user data
         */
        function removePopup(data) {
            layui.layer.confirm(lang['ConfirmRemoveUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Remove, [data], index);
            });
        }

        /**
         * disable one user popup
         * @param data {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} user data
         */
        function disablePopup(data) {
            layui.layer.confirm(lang['ConfirmDisableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Disable, [data], index);
            });
        }

        /**
         * enable one user popup
         * @param data {{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}} user data
         */
        function enablePopup(data) {
            layui.layer.confirm(lang['ConfirmEnableUser'], {
                title: lang['OperationConfirm'],
                btn: [lang['Confirm'], lang['Cancel']]
            }, function (index) {
                operate(apiType.Enable, [data], index);
            });
        }

        /**
         * operate actions
         * @param type {apiType} action type
         * @param data {[{user:string, token:string, comment:string, status:boolean, ports:string, domains:string, subdomains:string}]} user data list
         * @param index popup index
         */
        function operate(type, data, index) {
            var url;
            var extendMessage = '';
            if (type === apiType.Remove) {
                url = "/remove";
                extendMessage = ', ' + lang['RemoveUser'] + lang['TakeTimeMakeEffective'];
            } else if (type === apiType.Disable) {
                url = "/disable";
                extendMessage = ', ' + lang['RemoveUser'] + lang['TakeTimeMakeEffective'];
            } else if (type === apiType.Enable) {
                url = "/enable";
            } else {
                layer.layer.msg(lang['OperateError']);
                return;
            }
            var loading = layui.layer.load();
            $.post({
                url: url,
                type: 'post',
                contentType: 'application/json',
                data: JSON.stringify({
                    users: data
                }),
                success: function (result) {
                    if (result.success) {
                        reloadTable();
                        layui.layer.close(index);
                        layui.layer.msg(lang['OperateSuccess'] + extendMessage, function (index) {
                            layui.layer.close(index);
                        });
                    } else {
                        errorMsg(result);
                    }
                },
                complete: function () {
                    layui.layer.close(loading);
                }
            });
        }

        /**
         * reload user table
         */
        function reloadTable() {
            var searchData = layui.form.val('searchForm');
            layui.table.reloadData('tokenTable', {
                where: searchData
            }, true)
        }

        /**
         * show error message popup
         * @param result
         */
        function errorMsg(result) {
            var reason = lang['Other Error'];
            if (result.code === 1)
                reason = lang['ParamError'];
            else if (result.code === 2)
                reason = lang['UserExist'];
            else if (result.code === 3)
                reason = lang['ParamError'];
            else if (result.code === 4)
                reason = lang['UserFormatError'];
            else if (result.code === 5)
                reason = lang['TokenFormatError'];
            layui.layer.msg(lang['OperateFailed'] + ',' + reason)
        }

        /**
         * click event
         */
        $(document).on('click.search', '#searchBtn', function () {
            reloadTable();
            return false;
        }).on('click.reset', '#resetBtn', function () {
            $('#searchForm')[0].reset();
            reloadTable();
            return false;
        });
    }

    return loadUserList;
})(layui.$)
