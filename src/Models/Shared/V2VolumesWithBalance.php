<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2VolumesWithBalance
{
    /**
     *
     * @var string $account
     */
    #[\JMS\Serializer\Annotation\SerializedName('account')]
    public string $account;

    /**
     *
     * @var string $asset
     */
    #[\JMS\Serializer\Annotation\SerializedName('asset')]
    public string $asset;

    /**
     *
     * @var int $balance
     */
    #[\JMS\Serializer\Annotation\SerializedName('balance')]
    public int $balance;

    /**
     *
     * @var int $input
     */
    #[\JMS\Serializer\Annotation\SerializedName('input')]
    public int $input;

    /**
     *
     * @var int $output
     */
    #[\JMS\Serializer\Annotation\SerializedName('output')]
    public int $output;

    /**
     * @param  ?string  $account
     * @param  ?string  $asset
     * @param  ?int  $balance
     * @param  ?int  $input
     * @param  ?int  $output
     */
    public function __construct(?string $account = null, ?string $asset = null, ?int $balance = null, ?int $input = null, ?int $output = null)
    {
        $this->account = $account;
        $this->asset = $asset;
        $this->balance = $balance;
        $this->input = $input;
        $this->output = $output;
    }
}